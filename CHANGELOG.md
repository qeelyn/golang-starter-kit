# golang starter kit Change Log

本更新涉及的包更新记录包括go-common,gin-contrib

### 2018-09-28

* enh: http,RPC通讯增加gzip支持

### 2018-09-28
* enh: 增强graphql的错误定制输出
* enh: jwt验证完善,与go-common组件统一
* fix: auth验证链的传递
* fix: 用户ID 由 userId统一至小写 userid

### 2018-09-23
* enh: go-common为grpx服务定义增加grpc原生ServerOption支持
* enh: 客户端rpc配置,及服务端grpc服务配置定义
* enh: 日志记录整合grpc日志,请区分gateway及grpc服务端的调用方式.
* enh: go-common logger优化,增加上下文信息记录支持,区分类型化日志及糖方式 v