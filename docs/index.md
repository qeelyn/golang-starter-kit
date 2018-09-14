# 开发手册

## 项目目录
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
        - gateway.yaml          应用配置文件
        - gateway-local.yaml    个人的环境配置文件
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

基于数据库的应用在配置完后,就可以进行服务响应测试了.更多的配置内容可查看[应用程序配置](./application.md)

[下一节 应用结构](./application.md)