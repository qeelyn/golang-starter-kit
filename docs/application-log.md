# 日志

本套件的日志以uber.zap为核心,具有一定的定制性及可扩展性,可以很容易的记录各种类型的消息并分类处理,并将他们收集至特定的存储位置.

默认提供access.log及app.log,来记录访问日志与应用日志.

## 日志消息

跟系统日志一样,提供以下方法,
  - Strict 该方法提供类型化日志
  - Sugared 简单方式,与原生log方法原型最接近.
  - WithContext 会把上下文相关信息增加到日志字段中. 

配置
```
log:
  file:
    filename: "runtime/app.log"
    maxsize: 500
    maxbackups: 3
    maxage: 28
    level: -1
  access:
    filename: "runtime/access.log"
    maxsize: 200
    maxbackups: 3
    maxage: 28
    level: 0
```
目前的日志主要还是采用文件日志,后续再考虑结合其他方式,如Logstash

在api的访问日志中采用了zap记录,可根据需求自行更改显示的格式
```
    logger.Info(path,
        zap.Int("status", c.Writer.Status()),
        zap.String("method", c.Request.Method),
        zap.String("path", path),
        zap.String("query", query),
        zap.ByteString("body", bodyCopy.Bytes()),
        zap.String("ip", c.ClientIP()),
        zap.String("auth", c.GetString("userid")),
        zap.String("user-agent", c.Request.UserAgent()),
        zap.String("time", end.Format(timeFormat)),
        zap.Duration("latency", latency),
    )
```
## 使用方式

gateway:
```
app.Config.Strict().Error(...)  //zap方式
app.Config.WithContext().Error(....) //附带上下文的zap方式
app.Config.Sugared().Error(....) //
```
micro:

在微服务下,采用拦截器方式提供上下文信息的记录支持,大部分情况你应该使用如下方式进行记录
```
ctxzap.Extract(ctx).Error("query holding error:" + err.Error())
```
> 取消了之前在微服务构建增加Log组件的方法.统一使用ctxzap

[下一节 Service](service-layer.md)