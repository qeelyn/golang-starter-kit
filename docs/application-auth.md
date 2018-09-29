# 认证与授权

## 认证

本starter kit采用Json Web Token来支持常见系统的认证及授权.涉及到认证服务系统在此不多解释.

未来有机会也提供认证服务系统

经由用户认证服务获取的JWT在系统中通过各种方式传递.
* gateway

gateway接受http请求,接收Http Header的Authorization头信息传递.在Log Middleware中进行将该值记录到上下文向后传递.
```
func AccessLogHandleFunc(logger *zap.Logger, timeFormat string, utc bool) gin.HandlerFunc {
    if orgId := c.GetHeader("Qeelyn-Org-Id"); orgId != "" {
    	c.Set("orgid", orgId)
    }
    // pass to context
    if authHeader := c.GetHeader("Authorization"); authHeader != "" {
    	c.Set("authorization", authHeader)
    }
    c.Next()
}
```
* Gateway向GRPC服务传递

利用grpc client加载Auth拦截器以获取上下文中的Header信息
```
    cc, err := dialer.Dial(viper.GetString("name"),
		dialer.WithUnaryClientInterceptor(
			authfg.WithAuthClient(isGateway),
		),
	)
```
* GRPC服务接收

利用grpc_auth中间件接收MD对象中的Header信息
```
    var opts = []grpcx.Option{}
    opts = append(opts, grpcx.WithAuthFunc(
        authfg.ServerJwtAuthFunc(viper.GetStringMap("jwt"))
        )
    )
```
## 授权

* mvc
* graphql

由于GraphQl不像传统的MVC那样具有规则,经常做法就是独立设置授权编码,并检验.
可见gateway的CheckAccess组件实现.
(待细化)