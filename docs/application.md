应用结构
==============

应用
--------
应用可定义为单一服务程序的管理主体,每个服务程序只能包含一个应用主体,应用主体是通过app包进行访问

### 应用主体配置
如下所示,应用主体的配置目前还比较简单
```
// 应用程序名
appname: "myApp"
// 监听地址 host:port
listen: ":9097"
// 应用模式,生产或测试环境切换
appmode: debug
```
组件级别在应用配置的一级节点展开

应用主体是个很灵活的,应该而由项目管理者自行决定引入的组件,这边是有一些代码成本,但实际上带来的自由度让人感觉更好.

### 应用组件

以下列出一些顶层组件,

* [缓存](application-cache.md)
* [数据库](application-db.md)
* [授权](auth.md)
* [OpenTracing](application-opentracing.md) 分布式跟踪
* [metrics](application-metrics.md) 系统运行指标监控

[下一节 模型定义-protobuf](use-protobuf.md)