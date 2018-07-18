基于Go的应用开发入门套件
========================

本工具包旨于让您快速构建起项目结构,以便通过Go来开发WebApi或RPC服务,遵循SOLID的最佳实践来编写GO代码

本工具包提供下列功能:

* 应用与组件的可配置性
* 基于Gin的Web服务支持
* GraqhQl服务支持
* 基于Gorm的数据库操作及事务控制
* JWT-based 验证
* 异常处理及可控的错误响应
* 应用日志及访问日志支持
* 围绕protobuf为模型中心,生成通用性代码
* 采用Service层,并可扩展为RPC服务或微服务
* 测试环境可配置

本工具包使用了常见的GoPKG,你可以很容易的替换为自己喜欢的包.因为这些流行的PKG进行了良好的抽像.

* 路由框架: [gin](http://github.com/gin-gonic/gin)
* 数据库及ORM: [gorm](http://github.com/jinzhu/gorm)
* 数据验证: 目前通过Gin在路由层处理,还有很式工作 [want help]
* 配置文件: [viper](http://github.com/spf13/viper)
* 日志: [Uber Zap](http://go.uber.org/zap)
* graphql: [gopher-graphql](github.com/graph-gophers/graphql-go)
* 依赖管理: [DEP](https://golang.github.io/dep/docs/introduction.html)
* 基础套件:[qeelyn-common](http://github.com/qeelyn/go-common)
  - 缓存 cache 内置支持local,redis,memcached
  - protobuf工具包
  - grpc 一些的微服务工具包
* 中间件与组件: [qeelyn-contrib](http://github.com/qeelyn/gin-contrib)
* protoc生成工具扩展: [protoc-gen-goql](http://github.com/tsingsun/protoc-gen-goql)

微服务

* 服务注册与发现: 实现了[etcd](https://github.com/coreos/etcd),留有其他组件扩展的能力
* GRPC组件: 主要采用了[grpc-ecosystem](https://github.com/grpc-ecosystem)提供的组件
* 系统监控: [prometheus](https://prometheus.io),可配合[grafana]()https://grafana.com)搭建监控平台

本套件可以做什么
----------------

本套件面向是的企业级应用开发,做为通用的API编程框架.包括常见的RESTapi,微服务架构支持.

本套件的目标不是为了实现像beego这样的全栈框架,通常认为每个项目特性不同,除了提供一些基础包,应该由项目自行装配.

开发环境
---------

- go环境安装
- IDE vscode or goland
- go的开发离开不了翻墙,具体可看[工具篇](./docs/tools.md)  

快速入门
---------
### 新建项目
```
  git clone https://github.com/qeelyn/golang-starter-kit.git $GOPATH/src/xxx.com/your-vendor/project-name
```

项目目录
----------
### 项目结构

```
- gateway                   api gateway接口项目,基于Web http server
    - app                   应用程序域相关组件,包含配置,日志及启动有关的处理
    - config                配置文件目录
    - controllers           控制器文件夹
    - public                站点静态文件目录
    - routers               路由配置目录
    - schema                graphql定义文件夹
    - resolver              graphql golang解释器目录
    - loader                graphql 的dataloader目录
    server.go               应用服务文件
- cmd                       程序运行及配置
    - config
        - app.yaml          应用配置文件
        - app-local.yaml    个人的环境配置文件
    - srv                   服务目录,主要针对Server的初始化
        - gateway
    main.go                 入口,具体由决定
- schemas                   protobuf协议定义目录
- services                  逻辑服务代码文件夹
    - greeter.go            RPC服务的Service文件,包括应用程序域组件
```
配置
---------

在config目录中新建本地配置文件,有命名要求: app-local.yaml

默认所有配置应该集中于app.yaml与app-local.yaml中

* 应用配置
```
appname: myApp
listen: ":9097"
appmode: debug
web:
  staticdir: public
  
```

* 日志配置
```
log:
  file:
    filename: runtime/app.log"
    maxsize: 500
    maxbackups: 3
    maxage: 28
    level: -1
  access:
    filename: runtime/access.log"
    maxsize: 200
    maxbackups: 3
    maxage: 28
    level: 0

```
* DB配置 - 具体数据库类型请查看grom的配置
```
db:
  default:
    dialect: mysql
    dsn: root:@tcp(localhost:3306)/yak
```

基于数据库的应用在配置完后,就可以进行服务响应测试了.更多的配置内容可查看[应用程序配置](./docs/application.md)

[下一节 应用结构](./docs/application.md)