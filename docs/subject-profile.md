# 性能分析

go提供了基础的分析工具pprof,下面介绍一些常用的

* qcachegrind: 通过GUI查看pprof导出的报告
* wrk: 测试请求工具

## 开始

### 嵌入pprof代码
```
import _ "net/http/pprof"
```
### 运行HTTP服务器
```
go func() {
   http.ListenAndServe("localhost:8081", nil)
}()
```

目标服务启动后

* 请求端启动
```
wrk -c 200 -t 4 -d 3m -s pprof.lua http://localhost:8040/v2/query
```


### 使用pprof

服务端用以下命令进入跟踪
```
go tool pprof -seconds 200  ./all_go http://localhost:8081//debug/pprof/profile
```
在指定时间结束后,将进入pprof命令行:
* 导出qcachegrind可识别文件
```
callgrind
```

