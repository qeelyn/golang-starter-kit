日志
========
本套件的日志以uber.zap为核心,具有一定的定制性及可扩展性,可以很容易的记录各种类型的消息并分类处理,并将他们收集至特定的存储位置.

默认提供access.log及app.log,来记录访问日志与应用日志.

日志消息
--------------
跟系统日志一样,提供以下方法,
  - Debug,Debugf
  - Info,Infof
  - Warn,Warnf
  - Error,Errorf
  - Print,Printf

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
目前的日志主要还是采用文件日志,后续再考虑结合其他Logstash这样的,一劳永逸.

日志只是app包中的一部分,根据需求自行更改显示的格式
```
logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.ByteString("body", bodyCopy.Bytes()),
			zap.String("ip", c.ClientIP()),
			zap.String("auth", c.GetString("userId")),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", end.Format(timeFormat)),
			zap.Duration("latency", latency),
		)
```

[下一节 Service](service-layer.md)