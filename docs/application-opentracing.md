OpenTracing 与 Jaeger
======================

可通过[OpenTracing](./subject-opentracing.md)了解基本的知识

配置
-----------
```
opentracing:
  serviceName:"" // 该配置不需要的,自动采用appname
  sampler:
    type: const
    param: 1
  reporter:
    logSpans: false
    localAgentHostPort: "127.0.0.1:6831"
```
主要配置组件为
* sampler: 取样
* reporter: 提交

详细文档https://github.com/jaegertracing/jaeger-client-go

只要提供配置后,就可以在jaeger服务端查看到相应的查询信息

部署
------------
测试环境可直接用 jaegertracing/all-in-one:latest 进行部署
```
docker pull jaegertracing/all-in-one:latest
```

实际生产中.agent需要独立部署,以处理应用产生的span数据