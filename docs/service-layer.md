Service
==============
为了兼容单体服务与微服务架构,业务逻辑采用Repository模式,但为直接提供Service.

可以认为由protobuf生成的文件包含了IRepostiory定义,Service只是对这些契约的实现.

Service开发时,应该注意,这是可扩展为独立服务程序的,在组件应用时应该特别注意,如不要引用WebSite的包.特别是Gin的作用范围内的包.
在注意这边细节后,开发一个service主要就是纯业务逻辑的实现

为了做为可分离共用的服务层,应用程序组件App是独立的,具有独立的初始化步骤.
```go
type HelloService struct{
    //可以定义一些依赖组件,目前的组件直接使用App包,看个人喜欢进行调整
}

func NewHelloService(){
    return &HelloService{}
}

func World(ctx context.Context,req *WorldRequest) (*WorldResponse,error){
    db = app.Db
    // ...业务逻辑实现
    return &WorldResponse{Data:"hello world"},nil
}
```
调用
--------------
如果是单体项目,可在handle中这样定义,直接引入service层.如果需要单例管理,可以再多定义如ClientManager这样的管理类来维持.
而做为微服务提供的服务层是个独立的应用程序,采用GRPC服务.

### 调用 - 单体程序
```go
type HelloController struct {
	client     *hello.HelloService	
}

func ServeHelloResource(group *gin.RouterGroup) {
	helloCtr := HelloController{
		client:     NewHelloService(),		
	}
	group.GET("/hello", helloCtr.world)	
}
func (t *HelloController)world(c *gin.Context){
    //todo
}
```

GRPC
---------
做为微服务支持,将Service变成服务代码时,如果遵循Service的开发约定,只需要增加服务应用初始化即时

请直接转[微服务专题](./subject-micro.md)

[下一节 配合数据库工作](application-db.md)