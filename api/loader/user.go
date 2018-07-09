package loader

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
	"github.com/qeelyn/golang-starter-kit/api/app"
	"github.com/qeelyn/golang-starter-kit/helper/relay"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type UserLoader struct {
	userInfoUrl string
}

func NewUserLoader(opts ...dataloader.Option) *dataloader.Loader {
	val := UserLoader{
		userInfoUrl: app.Config.GetString("auth.auth-server") + "/userinfo",
	}.loadBatch
	return dataloader.NewBatchedLoader(val, opts...)
}

func (t UserLoader) loadBatch(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(keys)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)
	wg.Add(n)
	for i, key := range keys {
		go func(i int, key dataloader.Key) {
			defer wg.Done()
			var oid int
			if err := relay.UnmarshalSpec(graphql.ID(key.String()), &oid); err != nil {
				results[i] = &dataloader.Result{Data: nil, Error: err}
				return
			}
			gin := ctx.Value(0).(*http.Request)
			if userInfo, err := t.getUserInfo(gin, strconv.Itoa(oid)); err != nil {
				results[i] = &dataloader.Result{Data: nil, Error: err}
			} else {
				results[i] = &dataloader.Result{Data: userInfo, Error: err}
			}
		}(i, key)

	}
	wg.Wait()
	return results
}

func (t UserLoader) getUserInfo(context *http.Request, openId string) (map[string]interface{}, error) {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("GET", t.userInfoUrl+"?openid="+openId, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", context.Header.Get("Authorization"))
	authRes, err := client.Do(req)
	if err == nil {
		var body []byte
		body, err = ioutil.ReadAll(authRes.Body)
		if err != nil {
			return nil, err
		}
		var res map[string]interface{}
		json.Unmarshal(body, &res)
		if authRes.StatusCode != http.StatusOK {
			msg := res["errors"].([]interface{})[0].(map[string]interface{})["message"]
			err = errors.New(msg.(string))
		} else {
			return res, nil
		}
	}
	err = fmt.Errorf("error on auth client request : %s", err)
	return nil, err
}

func LoadUserNickName(ctx context.Context, key dataloader.Key) (string, error) {
	if val, err := Load(UserLoaderKey, ctx, key); err != nil {
		return "", err
	} else {
		ui := val.(map[string]interface{})
		return ui["data"].(map[string]interface{})["nickname"].(string), nil
	}
}
