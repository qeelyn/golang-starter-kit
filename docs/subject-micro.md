# 微服务-microservices
本开发套件提供了很简单的方式来支持微服务架构.具备以下能力

* 服务发现与注册: 采用基于etcd的方式,支持`HOSTIP`环境变量,简化配置
* 负载均衡: GRPC提供,默认采用round_robin
* 消息格式: 目前为默认的GRPC的protobuf格式
* 消息流: 内置的拦截器支持unary和stream方式
* 跟踪服务: uber jeager
* 日志收集: uber zap,计划采用ELK技术来完善日志收集流程
* 监控与预警: prometheus + grafana

这些能力可通过grpcx.Option来初始化开启

## 服务构建

可通过`go-common/grpx`包提供的默认定义初始化.用户定制要求高的,可参考mirco方法实现
```    
    tracer := tracing.NewTracer(app.Config.Sub("opentracing"), appName)
	serverPayloadLoggingDecider := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
		if fullMethodName == "healthcheck" {
			return false
		}
		return true
	}
    // option 初始化
	var opts = []grpcx.Option{
		grpcx.WithTracer(tracer),
		grpcx.WithLogger(app.Logger.GetZap()),
		grpcx.WithUnaryServerInterceptor(
		    grpc_zap.PayloadUnaryServerInterceptor(
		        app.Logger.GetZap(),
		        serverPayloadLoggingDecider)),
		grpcx.WithAuthFunc(grpcx.AuthFunc(config.GetString("auth.public-key"))),
		grpcx.WithPrometheus(config.GetString("metrics.listen")),
		grpcx.WithRegistry(register, fofSrvName),
	}
    //采用默认的微服务组件初始化
	server, err := grpcx.Micro(appName, opts...)

	if err != nil {
		panic(fmt.Errorf("fof server start error:%s", err))
	}
    
	rpc := server.BuildGrpcServer()
	// grpc的服务注册
	fof.RegisterFofServiceServer(rpc, fofsrv.NewFofService())
	// 启动
	if err = server.Run(rpc, listen); err != nil {
		return fmt.Errorf("Server run error:", err)
	}
	return nil
```

### 注册配置

以qeelyn://author为基地址,相关参数如下
```
appname: srv-notice  //注册中心路径为: qeelyn://author/srv-notice
registryListen: ":8033" 服务地址
```

如果设置了HOSTIP环境变理,配置文件为`:8033`格式的变根据环境变量修改成`${HOSTIP}:8033`
该支持主要为让不同机器保持相同的配置.

### go客户端

基于grpc的注册与发现可通过配置的方式
* 服务名: 采用GRPC的naming方式,格式如: qeelyn://author/srv-pool
* IP: 传统方式,如"127.0.0.1:8000"
```
func newDialer(serviceName string, tracer opentracing.Tracer) *grpc.ClientConn {
	cc, err := dialer.Dial(serviceName,
		dialer.WithDialOption(
			grpc.WithInsecure(),
			grpc.WithBalancerName("round_robin"),
		),
		dialer.WithUnaryClientInterceptor(
			grpc_prometheus.UnaryClientInterceptor,
			dialer.WithAuth(),
		),
		dialer.WithTracer(tracer),
	)
	if err != nil {
		log.Panicf("dialer error: %v", err)
	}
	return cc
}
    //服务名
	poolcc := newDialer(app.Config.GetString("rpc.pool"), tracer)	
    //grpc client
	app.PoolClient = pool.NewPoolServiceClient(poolcc)

```

### 复杂网络下的rpc通信

* 心跳
采用gRPC长连接方式构建微服务通信,已经考虑到心跳设计,对于简单环境下,默认即可.但在复杂网络环境下,必须考虑无活动连接被防火墙释放问题.如
客户IDC网络与云环境的联通情况下,需要自定义心跳,特别需要采用双向心跳,一般只需要定义发送间隔.

    客户端: 在rpc下定义keepalive
    服务端: 在根配置下定义keepalive
    
* 断线重连
由gRPC定义实现
