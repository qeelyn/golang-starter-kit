# 日志

这些我们先介绍elastic全家桶

## 采集

## filebeat

### docker使用

* 获取镜像

使用elastic自家的镜像仓库
```
docker pull docker.elastic.co/beats/filebeat:6.4.2
```

* 运行
```
docker run \
  --mount type=bind,source="$(pwd)"/filebeat.yml,target=/usr/share/filebeat/filebeat.yml \
  docker.elastic.co/beats/filebeat:6.4.2
```

