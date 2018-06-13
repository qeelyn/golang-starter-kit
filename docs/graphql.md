Graphql
==============
除了采用传统的WEB API方式外,还可以采用graphql的方式提供,实际上,这种方式已经越来越被大家接受,相比web api的优势明显.
* 字段级别的解释,变更可控性强.减少很多烦人的版本控制
* 统一的编程模型,方便跟踪与性能提升,比api这样的黑洞更友好.
* 基于模型的定义访问,能应对各种需求变化,原型上,基于同一套模型的前端设计,接口是不需要更改的.

[graphql-go]((github.com/graph-gophers/graphql-go)的参考资料较少,本质上需要对js版本的graphql有一定的了解,再参考示例代码,学习一下基础知识.

graphql属于Web层,与gin是相配置的.

模型定义
-----------
采用JS同样的定义支持,采用golang的资源加载方式,目录结构如下
```
api
  - schema //在该文件夹存放graphql定义相关文件
    - type  //类型文件夹 在该文件夹定义类型
    - schema.graphql
    - schema.go  //资源加载,graphql只需要将内容进行拼接
```
Resolver
-------------
对定义文件的解释器,为golang实现代码

Tracer
-------------
由于graphql的灵活性,可能一个类型中,就会并发访问不同的remote service,这样在跟踪及性能分析上就存在一定的困难,需要借助一些分析工具.
目前采用 opentracing + jaeger 来提供跟踪服务

后续再来专门讲.