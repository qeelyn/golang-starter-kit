处理请求
==============

运行机制
--------------

经过应用主体的路由机制后,采用的是洋葱模型,由一层层中间件来完成整个请求的处理,具体请参考Gin的处理方式,

路由
--------------
路由的初始化文件在router文件夹中.一般在本router.go文件中完成路由分组与中间件的集成


可直接参考Gin的方式来编写web服务,如果你钟情于MVC,就自定义controller文件夹.本套件建议就用GO的方式定义API,即直接定义Handle方法

```
	v1 := g.Group("oauth2/v1")
	{
		v1.POST("token", tokenHandle)
		v1.POST("authorize", authorizeHandle)
	}

	authMd := app.BearerAuth(app.Config.GetStringMap("auth"))

	g.POST("login", api.Login)
	g.POST("logout", authMd, api.Logout)
	g.GET("userinfo", authMd, api.UserInfo)

```

请求参数
-------------
Gin已经提供了方便的处理方法,我们可以直接使用它
```
func Get(c *gin.Context){
    c.Qeury("key")
    c.Param("key")
}
```
### 参数绑定
```
gin.Context.Bind(struct)
```

响应
----------

除了Gin提供的方法外.你完全可以定义一些helper方法来辅助格式化响应
```
gin.Json
gin.Html
等等
```

### gin的作用范围

gin框架的作用范围在路由及handle方法,gin.context的显式调用就结束,此时上下文会转化为context.Context接口或其封装往下传递.如service层或者graphql服务中传递,不推荐再将context转化为gin.context.
> 实际上将context转化为gin.context这也是很困难的,封装的层级是未知的.

[下一节 异常处理](error-handle.md)