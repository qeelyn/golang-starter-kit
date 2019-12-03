##  代理设置

在Go 1.13中，我们可以通过GOPROXY来控制代理，以及通过GOPRIVATE控制私有库不走代理。

设置GOPROXY代理：

    go env -w GOPROXY=https://goproxy.cn,direct
    设置GOPRIVATE来跳过私有库，比如常用的Gitlab或Gitee，中间使用逗号分隔：
    go env -w GOPRIVATE=*.gitlab.com,*.gitee.com

如果在运行go mod vendor时，提示Get https://sum.golang.org/lookup/xxxxxx: dial tcp 216.58.200.49:443: i/o timeout，则是因为Go 1.13设置了默认的GOSUMDB=sum.golang.org，这个网站是被墙了的，用于验证包的有效性，可以通过如下命令关闭：

    go env -w GOSUMDB=off
 
可以设置 GOSUMDB="sum.golang.google.cn"， 这个是专门为国内提供的sum 验证服务。

    go env -w GOSUMDB="sum.golang.google.cn"

> 有些IDE(如goland)会集成包管理工具,需要注意一下,是否存在环境变量被覆盖的情况

## go < 1.13
### 翻墙

目前套件使用翻墙的主要原因为依赖golang.org的包.而golang.org是被墙的.

但发现在使用GOPATH的包时编译速度下降,建议有条件翻墙就上.

* 使用带有http proxy,如shadowsocks(mac),因为如果使用全局模式时,在socks模式下会被站点墙
   
   1. 设置好terminal: 
    ``` 
    // linux
    export https_proxy=http://127.0.0.1:1087
    
    // window(在win10、shadowsocks PAC模式下测试过,goland中配置proxy无效)
    set https_proxy=http://localhost:1080
    ```
### DEP(deprecated)

```
go get -u github.com/golang/dep/cmd/dep

dep ensure 同步包
```

第一次初始化时,可用以下命令
```
    dep init -gopath -v -no-examples                
```    

### VGO

在Go 1.11后,vgo正式成为go工具链的一部分,也意味着go官方正式推出版本管理工具.在笔者的使用过程来看,确实是最优秀的.

配合Gitee,完全可以不用翻墙,由于1.11还未发布,这边就不写了,请根据官方文档来

可以将被墙的包通过go.mod文件的replace方式替换,竟味着可以直接指向github.
如果github也被墙,可以通过gitee导入github的项目,然后将包指向gitee,这样就可以达到不用翻墙开发了
替换获取包常见的命令为
```
vgo get package@[commit|version]
```
> 通过gitee访问github,非常之快.目前为止,github又不抽疯了..请忽略gitee的包
```
//以下列出常见需要替换的包
replace (
    //本项目需要
    github.com/graph-gophers/graphql-go => github.com/qeelyn/graphql-go v0.0.0-20181012014650-03df3acf1181
	// github
	golang.org/x/net => github.com/golang/net v0.0.0-20180811021610-c39426892332
	golang.org/x/sys => github.com/golang/sys v0.0.0-20180810173357-98c5dad5d1a0
	golang.org/x/text => github.com/golang/text v0.3.0
	google.golang.org/appengine => github.com/golang/appengine v1.1.0
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20180808183934-383e8b2c3b9e
	google.golang.org/grpc => github.com/grpc/grpc-go v1.14.0
    
    // gitee.com/githubmirror.
	golang.org/x/net => gitee.com/githubmirror/golang-net v0.0.0-20180811021610-c39426892332
	google.golang.org/appengine => gitee.com/githubmirror/appengine v1.1.0
	google.golang.org/genproto => gitee.com/githubmirror/go-genproto v0.0.0-20180808183934-383e8b2c3b9e
	google.golang.org/grpc => gitee.com/githubmirror/grpc-go v1.14.0
)
```

> 现发现vgo的问题为,加载不出非go文件目录.对于protobuf需要外部protobuf文件编译会产生问题,请自行COPY

### Build

本套件默认的编译路径为cmd,所有的配置也是针对该路径的.请在开发时调整一下IDE配置
* goland
![img](./img/goland-build.jpg)