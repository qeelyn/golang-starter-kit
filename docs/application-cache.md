缓存
===========
通过对比现在go的一些应用缓存组件,比如beego的cache组件,让人意外的是这些组件的能力太弱,如Get返回Interface,支持的数据类型有限.只好自行实现.
提供如下功能.
* 提取自动根据入参识别并转换类型.
* 全类型的支持 保存至缓存.
* 采用msgpack进行序列化.
* 内置支持local,redis,memcache
* 多缓存组件并存,支持组件命名

配置
--------
```
cache:
  default:
    type: memcache
    addr: :11211
  auth:
    type: redis
    addr: :6379
    password:
    db: 3
    connectionTimeout: 3
  local:
    type: local
    duration: 10 #默认时间 分钟
    gc: 30 # 分钟
```