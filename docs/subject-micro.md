微服务-microservices
====================
本开发套件提供了很简单的方式来支持微服务架构.具备以下能力

* 服务发现与注册: 采用基于etcd的方式
* 负载均衡: GRPC提供,默认采用round_robin
* 消息格式: 目前为默认的GRPC的protobuf格式
* 消息流: 内置的拦截器支持unary和stream方式
* 跟踪服务: uber jeager
* 日志收集: uber zap,计划采用ELK技术来完善日志收集流程
* 监控与预警: prometheus + grafana

这些能力可通过grpcx.Option来初始化开启

服务构建
-------------
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

go客户端
----------

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