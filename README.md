基于Go的应用开发入门套件
========================

本工具包旨于让您快速构建起项目结构,以便通过Go来开发WebApi或RPC服务,遵循SOLID的最佳实践来编写GO代码.

> 现阶段本项目例子还相对简单,但框架及组件的使用都是在实际项目使用的,未来再提供近于实战的例子.

本工具包提供下列功能:

* 应用与组件的可配置性,并支持配置中心方式
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
* 依赖管理: [DEP](https://golang.github.io/dep/docs/introduction.html) 将被[vgo](https://github.com/golang/vgo)取代
* 基础套件:[qeelyn-common](http://github.com/qeelyn/go-common)
  - 缓存 cache 内置支持local,redis,memcached
  - protobuf工具包
  - grpc 一些的微服务工具包
* 中间件与组件: [qeelyn-contrib](http://github.com/qeelyn/gin-contrib)
* protoc生成工具扩展: [protoc-gen-goql](http://github.com/tsingsun/protoc-gen-goql)

微服务部分

* 服务注册与发现: 实现了[etcd](https://github.com/coreos/etcd),留有其他组件扩展的能力
* GRPC组件: 主要采用了[grpc-ecosystem](https://github.com/grpc-ecosystem)提供的组件
* 系统监控: [prometheus](https://prometheus.io),可配合[grafana]()https://grafana.com)搭建监控平台
* 极方面的通过Docker构建部署.可通过[基于jenkins的持续构建](./docs/subject-jenkins.md)进一部了解

本套件可以做什么
----------------

本套件面向是的企业级应用开发,做为通用的API编程框架.包括常见的RESTapi,微服务架构支持.

本套件的目标不是为了实现像beego这样的全栈框架,通常认为每个项目特性不同,除了提供一些基础包,应该由项目自行装配.

开发环境
---------

- go环境安装
- IDE vscode or goland
- 以前go的开发离开不了翻墙,现在可以不翻了,具体可看[工具篇](./docs/tools.md)  

快速入门
---------
### 新建项目
```
  git clone https://github.com/qeelyn/golang-starter-kit.git $GOPATH/src/xxx.com/your-vendor/project-name
```

### 运行
```
  make run    
```
看到启用界面时,通过浏览器访问: `http://localhost:18000/graphiql`,这时你看到graphiql都正常时.项目成功运行.

默认默认采用微服务划分,在graphiql中输入
```
query test {
  hello(name:"qsli") {
    id
  }
}
```
执行即可

### 资源

开发手册: [戳我](./docs/index.md) 

QQ: 21997272