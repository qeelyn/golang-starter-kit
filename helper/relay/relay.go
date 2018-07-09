package relay

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/graph-gophers/graphql-go"
	"strings"
)

// return type ,error
func ParseGlobalId(gid string, v interface{}) (string, error) {
	s, err := base64.URLEncoding.DecodeString(string(gid))
	if err != nil {
		return "", err
	}
	i := strings.IndexByte(string(s), ':')
	if i == -1 {
		return "", err
	}
	json.Unmarshal([]byte(s[i+1:]), v)
	return string(s[:i]), nil
}

func MarshalID(kind string, spec interface{}) graphql.ID {
	d, err := json.Marshal(spec)
	if err != nil {
		panic(fmt.Errorf("relay.MarshalID: %s", err))
	}
	return graphql.ID(base64.URLEncoding.EncodeToString(append([]byte(kind+":"), d...)))
}

func UnmarshalKind(id graphql.ID) string {
	s, err := base64.URLEncoding.DecodeString(string(id))
	if err != nil {
		return ""
	}
	i := strings.IndexByte(string(s), ':')
	if i == -1 {
		return ""
	}
	return string(s[:i])
}

func UnmarshalSpec(id graphql.ID, v interface{}) error {
	s, err := base64.URLEncoding.DecodeString(string(id))
	if err != nil {
		return err
	}
	i := strings.IndexByte(string(s), ':')
	if i == -1 {
		return errors.New("invalid graphql.ID")
	}
	return json.Unmarshal([]byte(s[i+1:]), v)
}
