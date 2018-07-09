OpenTracing
======================

[开放分布式追踪（OpenTracing）入门与 Jaeger 实现](https://yq.aliyun.com/articles/514488)

分布式系统的运维挑战
------------------
容器、Serverless 编程方式的诞生极大提升了软件交付与部署的效率。在架构的演化过程中，可以看到两个变化：

![](http://ata2-img.cn-hangzhou.img-pub.aliyun-inc.com/9e95e41a2fa724dc66b2f49a967845c6.png)
* 应用架构开始从单体系统逐步转变为微服务，其中的业务逻辑随之而来就会变成微服务之间的调用与请求
* 资源角度来看，传统服务器这个物理单位也逐渐淡化，变成了看不见摸不到的虚拟资源模式

从以上两个变化可以看到这种弹性、标准化的架构背后，原先运维与诊断的需求也变得越来越复杂。为了应对这种变化趋势，诞生一系列面向 DevOps 的诊断与分析系统，包括集中式日志系统（Logging），集中式度量系统（Metrics）和分布式追踪系统（Tracing）

### Logging，Metrics 和 Tracing

Logging，Metrics 和 Tracing 有各自专注的部分。

Logging - 用于记录离散的事件。例如，应用程序的调试信息或错误信息。它是我们诊断问题的依据。
Metrics - 用于记录可聚合的数据。例如，队列的当前深度可被定义为一个度量值，在元素入队或出队时被更新；HTTP 请求个数可被定义为一个计数器，新请求到来时进行累加。
Tracing - 用于记录请求范围内的信息。例如，一次远程方法调用的执行过程和耗时。它是我们排查系统性能问题的利器。
这三者也有相互重叠的部分，如下图所示:
![](http://ata2-img.cn-hangzhou.img-pub.aliyun-inc.com/92806aa2426813a4f47e6ba9b01f76f7.png)

通过上述信息，我们可以对已有系统进行分类。例如，Zipkin 专注于 tracing 领域；Prometheus 开始专注于 metrics，随着时间推移可能会集成更多的 tracing 功能，但不太可能深入 logging 领域； ELK，阿里云日志服务这样的系统开始专注于 logging 领域，但同时也不断地集成其他领域的特性到系统中来，正向上图中的圆心靠近。

关于三者关系的更详细信息可参考 [Metrics, tracing, and logging](http://peter.bourgon.org/blog/2017/02/21/metrics-tracing-and-logging.html?spm=a2c4e.11153940.blogcont514488.18.11b730c2dYH4KD)下面我们重点介绍下 tracing

### Tracing 的诞生
Tracing 是在90年代就已出现的技术。但真正让该领域流行起来的还是源于 Google 的一篇论文"Dapper, a Large-Scale Distributed Systems Tracing Infrastructure"，而另一篇论文"Uncertainty in Aggregate Estimates from Sampled Distributed Traces"中则包含关于采样的更详细分析。论文发表后一批优秀的 Tracing 软件孕育而生，比较流行的有：

* Dapper(Google) : 各 tracer 的基础
* StackDriver Trace (Google)
* Zipkin(twitter)
* Appdash(golang)
* 鹰眼(taobao)
* 谛听(盘古，阿里云云产品使用的Trace系统)
* 云图(蚂蚁Trace系统)
* sTrace(神马)
* X-ray(aws)
分布式追踪系统发展很快，种类繁多，但核心步骤一般有三个：代码埋点，数据存储、查询展示。

下图是一个分布式调用的例子，客户端发起请求，请求首先到达负载均衡器，接着经过认证服务，计费服务，然后请求资源，最后返回结果。

![opentracing1.png](http://ata2-img.cn-hangzhou.img-pub.aliyun-inc.com/9e30ca9d4ac0f76502e1be2f87d2a2df.png)

数据被采集存储后，分布式追踪系统一般会选择使用包含时间轴的时序图来呈现这个 Trace。

![opentracing2.png](http://ata2-img.cn-hangzhou.img-pub.aliyun-inc.com/e8813c5c6d2bf7cdc4780ad7fe477245.png)

但在数据采集过程中，由于需要侵入用户代码，并且不同系统的 API 并不兼容，这就导致了如果您希望切换追踪系统，往往会带来较大改动。

其他内容请参考[原文](https://yq.aliyun.com/articles/514488)

### Jaeger

中文名应该是-耶格,由uber开发的实现OpenTracing API的系统.可通过Docker部署测试

阿里云提供了Jaeger的支持,其他云厂商还未查阅资料.