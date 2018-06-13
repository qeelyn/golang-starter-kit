Service
==============
为了兼容单体服务与微服务架构,业务逻辑采用Repository模式,但为直接提供Service.

可以认为由protobuf生成的文件包含了IRepostiory定义,Service只是对这些契约的实现.

Service开发时,应该注意,这是可扩展为独立服务程序的,在组件应用时应该特别注意,如不要引用WebSite的包.特别是Gin的作用范围内的包.
在注意这边细节后,开发一个service主要就是纯业务逻辑的实现
```
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

客户端调用
--------------
通过protobuf可生成各个语言的客户端调用代码,如果是单体项目,可在handle中这样定义,直接引入service层.如果需要单例管理,可以再多定义如ClientManager这样的管理类来维持
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
> RPC化目前项目还未开始实践,后续再补

[下一节 配合数据库工作](application-db.md)