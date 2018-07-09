使用protobuf
===============

protobuf是gRPC的协议,扩展名为<b>.proto</b>,利用一系列代码生成工具可以生成各层级代码.
[protobuf协议说明-要翻墙](https://developers.google.com/protocol-buffers/docs/reference/go-generated)

代码生成工具
------------

```
brew install protobuf
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u github.com/tsingsun/protoc-gen-goql
```

这里利用到了自定义的生成工具,主要作用是替换协议中默认的协议处理类,以使生成的代码可以在各逻辑层中传递.

* 处理像timestamp的非标量定义,控制序列化与反序列化
* 与GORM的配套使用,数据库命名推荐采用下划线方式,来匹配protobuf规范
* 支持采用inject注释为属性增加TAG定义,以及修改JSON定义
* nullable支持,但需要其他平台的客户端代码协同支持
* proto本身具备了json定义,可直接当成DTO
* request对象可直接做为gin的入参对象

通过proto定义输出的模型文件可用于数据处理层,传输层,省掉了我们重复定义模型及值传递这类无意义的过程

编译protobuf文件
------------------
```
protoc schemas/fund/fund.proto -I vendor/github.com/qeelyn/go-common -I . --goql_out=paths=source_relative:./
//需要rpc的话,加plugins=grpc
```

注意确保protoc-gen-go和protoc在PATH环境变量中，protoc-gen-go一般在$GOPATH/bin目录下。
> 注意paths参数,针对我们项目请加上该参数

* protbuf生成的文件没有的修改必要,如果需要修改,由生成器来调整.
* 请在protobuf文件头部注释出生成命令,执行命令的路径应该为当前项目目录.
* 推荐直接将RPC Service一起定义,后期可直接转RPC

### IDE支持

在IDE插件中，使用外部的定义时，可以配置导入相关的定义，如下图：
![pb-import](http://120.77.219.247/xhguo/fund_core_cn/uploads/3a6f08e3db4f6730d6d757dda97e2722/pb-import.png)

在Goland中，可以通过protobuf的插件来支持，在外部工具中配置编译器，如下图：
![pb-tool](http://120.77.219.247/xhguo/fund_core_cn/uploads/b65cc1b68f43c5e62ed84e58802279ca/pb-tool.png)
![pb-tool2](http://120.77.219.247/xhguo/fund_core_cn/uploads/d27867ce3d9bcc75ef00be29d6d1a60e/pg-tool-2.png)

[上一节 首页](README.md)    
  
[下一节 处理请求](handle-request.md)
