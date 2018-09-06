应用结构
==============

应用
--------
应用可定义为单一服务程序的管理主体,每个服务程序只能包含一个应用主体,应用主体是通过app包进行访问.

应用配置以文件形式展现,并支持远程配置.

为了便于在开发环境的个性化,支持本地配置文件(以-local为后缀的形式),并避免敏感信息的泄露.

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

### 本地配置

默认采用./config路径为应用配置路径,应用对应的配置文件已经在程序中指定完成.
当应用的配置文件为`gateway.yaml`时,将默认加载后缀为`-local`的配置文件`gateway-local.yaml`

需要注意当通过`-c`参数变更路径时,将与配置中心进行组合获取配置文件

### 远程配置
Kit默认以本地文件加载配置文件,如果启用远程配置时,注册中心与配置中心依赖服务是一致的,这时将同时启动服务注册与发现.

* 代码
```
	cnfOpts := config.ParseOptions(configOptions...)
	//采用etcd 客户端v3版本
	etcdv3.Build(cnfOpts)
	
	cnfOpts.FileName = "gateway.yaml"
	if app.Config, err = config.LoadConfig(&cnfOpts); err != nil {
    		return err
    }
```
* 启动参数
```
{cmd} -c {配置路径} -n {etcd connection string}
配置路径为uri : golang-start-kit/config 格式的基路径,
etcd connection string: 需要用host:port?key=value形式,如127.0.0.1:2379?username=qeelyn
例如启动gateway时,最终对应到etcd的配置路径为 : 127.0.0.1:2379/golang-start-kit/config/gateway.yaml,
所以往配置中心发布时,需要注意key取值
```
> etcd的连接串是根据etcd config的配置,当启用TLS时,需要定制初始化
* 发布配置

命令行方式可以用如下命令,自己可以结合一些开源的etcd管理工具或者自动开发,实现配置的管理
```
cat cmd/config/gateway.yaml | etcdctl put golang-start-kit/config/gateway.yaml
```

[下一节 模型定义-protobuf](use-protobuf.md)