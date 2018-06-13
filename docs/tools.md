### 翻墙

* 使用带有http proxy,如shadowsocks(mac),因为如果使用全局模式时,在socks模式下会被站点墙
   
   1. 设置好terminal: 
    ``` 
    // linux
    export http_proxy=http://127.0.0.1:1087
    export https_proxy=http://127.0.0.1:1087
    
    // window(在win10、shadowsocks PAC模式下测试过,goland中配置proxy无效)
    set http_proxy=http://localhost:1080
    set https_proxy=http://localhost:1080
    ```
### DEP

```
go get -u github.com/golang/dep/cmd/dep

dep ensure 同步包
```


第一次初始化时,可用以下命令
```
    dep init -gopath -v -no-examples                
```    
