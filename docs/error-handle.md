异常处理
==========
针对在响应过程,针对异常的响应是必不可少的环节.
golang自身的错误处理方式相对于其他语言可以说比较独特.异常做为方法的返回值处理.

Error只提供了如何转化为String方法,转化为文件输出.完全没有其他语言的序列化问题.这种单一性对于错误的提醒不够完备,因此需要自定义错误处理.

错误消息
-----------

套件中对错误消息的定义在config/errors.{语言}.yaml文件中.格式如下
```yaml
INTERNAL_SERVER_ERROR:
  code: 500
  message: "内部服务器错误!"
  debug: "Internal server error: {error}"

PERMISSION_DENIED:
  code: 403
  message: "无请求权限!"
  debug: "无请求权限: {error}"
```
INTERNAL_SERVER_ERROR 对应error的消息文本,通过消息文本的全等匹配来定位code与message,输出符合code,message这样的输出.

统一异常
----------

### gin框架

以gin为路由的错误处理是通过套件提供的ErrorHandle中间件实现的.异常的注入方式
```go
  gin.Context.Error(error)
```
最终ErrorHandle将error转化为code,message响应输出
> 涉及到gin的作用范围

### RPC

延用GRPC框架的异常处理.

### 异常跟踪

如果采用了OperatingTracer中间件时,会在Context中保存TracerID,来做为整个请求链的全局ID
```
// gin 
g *gin.Engine
g.Use(app.NewJeagerTracer())

// grpc client
cc, err := dialer.Dial(serviceName,		
	dialer.WithTracer(tracer),
)

// grpc server 初始化option时加入
var opts = []grpcx.Option{
    grpcx.WithTracer(tracer),
}
server, err := grpcx.Micro(appName, opts...)
```
访问日志及GRPC日志将会产生相应的key记录: "trace.traceid":"1fa3ff926212922"
同时tracerid可用于JeagerUI查询对应的operationtracer记录.

[下一节 日志](application-log.md)
