# 基于jenkins的持续集成

## 写在前面

本文内容把个性化的环境隐藏或替换,如xxx等.请按照自己情况调整

* docker源,选国内,如果是阿里云,请进入[加速器](https://cr.console.aliyun.com/#/accelerator),注意,本文写时,阿里上的配置是错误的,请下看
* docker的配置文件

  配置文件位于/etc/docker/daemon.json 没有的话自己建一个.
```
{
  "registry-mirrors": ["http://harbor.test.com"], #镜像加速地址
  "insecure-registries": ["xx.xx.xxx.231","registry.cn-shenzhen.aliyuncs.com"], # Docker如果需要从非SSL源管理镜像，这里加上。
  "max-concurrent-downloads": 10
}
```
>对于上传方也需要设置insecure-registries才可push至私有仓库,不过jenkins内置docker插件不用
## 安装

为了适用于各种开发环境,我们需要安装的以基础镜像打包的.注意不要采用基于alpine的,库支持不足,会遇到不少问题.
```
docker run \
  -u root \
  -d \
  -p 58080:8080 \
  -p 50000:50000 \
  -v jenkins-data:/var/jenkins_home \
  -v /usr/bin/docker:/usr/bin/docker \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --restart=always \
  jenkinsci/jenkins
```
### docker out docker

官方的安装方式只是将docker.sock映射内容器,但如果需要docker的话,会出现`docker not found`.因此你看到在run 命令中加入docker的映射.
但还是会出现一个libltdl.so.7文件不存在,这个可以进容器,执行apt-get安装
```
apt-get update && apt-get install -y libltdl7
```
> 可以自定义Jenkins的Dockfile,将so安装好.

### SSH KEY

进入容器,创建SSH KEY,这个KEY与GitLab会用到
```
ssh-keygen -t rsa -C "jenkins"
# 一路回车, 默认路径和文件名, 不要密码
cd ~/.ssh
//私钥和密钥可另保存起来,
mv id_rsa id_rsa_jk
mv id_rsa.pub id_rsa_jk.pub
```

在gitlab的项目下, 点击右侧配置菜单 -> Deploy Keys, 用刚才创建的 id_rsa_tho.pub 的内容, 创建一个key, 名称为 Readonly Key for Jenkins, 如果有多个项目都需要这个私钥, 则在每个项目的deploy keys下enable这个key即可. 可以用下面的方式验证是否生效

检验: 在Jenkins中, 新建一个freestyle的项目, 点击项目 -> Source Code Management, 选择 Git, 填入gitlab中给的项目地址 git@192.168.1.109:cc/tho.git 在下面add new credential, Username:git, Private Key Enter Directly, 输入刚才创建的 id_rsa_tho 的内容, 注意这个是私钥.

如果切换鼠标焦点后, 项目地址栏下方没有错误提示, 就说明Jenkins检查没错


## gitlab

重启
```
cd /data/gitlab-8.9.6
./ctlscript.sh restart
```
 
## registry

registry做为Docker镜像仓库,目前阿里云的容器镜像是免费的,可以优先使用.

[Docker官方文档](https://docs.docker.com/registry/deploying/#support-for-lets-encrypt)

registry官方的简易安装是不包含安全性的,我们直接采用auth方式

### auth
```
docker run \
  --entrypoint htpasswd \
  registry:2 -Bbn jenkins {pwd} > /opt/data/registry/auth/htpasswd

docker run -d \
  -p 5000:5000 \
  --restart=always \
  --name registry \
  -v /nas/registry:/var/lib/registry \
  -v /opt/data/registry/auth:/auth \
  -e "REGISTRY_AUTH=htpasswd" \
  -e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
  -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
  registry:2

```
检查: 
* netstat -an | grep 5000
* curl -x http://localhost:5000

### 客户端设置
insecure registries: 配置上私有的HOST:PORT,见/etc/docker/deamon.json

### 测试

找个包,tag 需要按registry的格式: HOST:PORT/{包包}
```
//如果配置了auth验证
docker login myregistrydomain.com:5000

docker tag busybox xx.xx.xxx.231:5000/busybox
docker push xx.xx.xxx.231:5000/busybox
```

## 基于pipeline的golang项目接入

### 设置GO基本环境

现在go为1.10版本,之前dep很流行,但很不幸,它无法解决翻墙问题.目前vgo是最适合国内的版本工具,而且将是go的标准工具.


* 默认的Jenkins镜像是不带有Go的编译工具的,安装GO Plugin插件.位置在/var/jenkins_home/tools下对应到插件位置.
* 进入镜像docker exec -it {containerid} bash
* ln -s /var/jenkins_home/tools/{}/bin/go /usr/bin/go
* 设置GOPATH为{$workspace}/go,建好{bin,pkg,src}三个目录
* 安装vgo:
```
// 自已将代码上传到GOPATH目录
go get golang.org/golang/x/vgo
ln -s /var/jenkins_home/workspace/go/bin/vgo /usr/bin/vgo
```
### 项目要求

采用vgo来做版本管理,并将一些无法获取的项目做replace

* 将包地址都改成国内码云的镜像地址,我使用githubmirror组同步github项目,但每天只能同像20个项目,珍惜吧.
* 所有以golang.org为结束的包

replace例子
```
replace(
    github.com/prometheus/client_golang => gitee.com/githubmirror/prometheus_client_golang v1
)
```


该包由于依赖特殊,还依赖github的一些项目,放入GOPATH会提示找不到相关包.
也是需要override的.

### pipeline例子
```
node {
    // Install the desired Go version
    def root = tool name: '1.10', type: 'go'
    def BUILD_NAME = "serve"
    def IMAGE_NAME = "xxx.xx.xx.244:5000/fof-api"
    def APP_ROOT = "src/github.com/tsingsun/fof-api"
    // Export environment variables pointing to the directory where Go was installed
    withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin","GOPATH=${WORKSPACE}/../go",]) {
        stage("prepare") {
            dir("src") {
                git branch: 'master',credentialsId: 'befbfc02-d877-4f6f-b3db-7ad4cb65119f', url: 'git@xx.xx.xxx.xx:advisorq/fof-api.git'
                sh "mkdir -p ${GOPATH}/src/github.com/tsingsun && rm -rf ${GOPATH}/${APP_ROOT} && ln -s ${WORKSPACE}/src ${GOPATH}/${APP_ROOT}"
            }
        }
        stage("build") {
            dir("src") {
                sh "vgo get && vgo mod vendor"
                sh "cp -rf gateway/public cmd/"
            }
            sh "cd ${GOPATH}/${APP_ROOT} && GOOS=linux GOARCH=amd64 go build -o cmd/${BUILD_NAME} cmd/main.go"
        }
        stage("deploy") {
            dir("src") {                
                docker.withRegistry('http://172.16.61.244:5000','pushToRegistry') {                    
                    def image = docker.build("${IMAGE_NAME}:${env.BUILD_ID}")
                    image.push("latest")    
                }
            }
        }
    }
}
```

## 部署

目前由于没有k8s与swarm环境,可采用dock pull的方式手工部署.
由于registry没有配套UI,但对于企业应用,它的restapi足够了.
```
//以下为安全验证下的请求
//列出image repo
http://jenkins:pwd@4xx.xx.xx.xxx:5000/v2/_catolog
//列出image tggs
http://jenkins:pwd@4xx.xx.xx.xxx:5000/v2/nginx/tags/list
```

 ## 常见问题
 
 * 忘记Jenkins管理员密码的解决办法
   按[此文可解](https://blog.csdn.net/jlminghui/article/details/54952148)
 * centos构建的无法在alpine运行.glibc库问题,但musl库是兼容的,可以用,在你的Dockfile加
 ```
 RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
 ```
 据说CGO编译也可以,但我没试
 
 * jenkins出现非登陆用户,这是根据git历史自动生成问题,这与git的用户配置有关.需要针对项目进行配置,在项目目录下
 ```
git config user.name “gitlab’s Name”
git config user.email “gitlab’s Name”
 ```
 
 ## 待完善
 * docker build的镜像是在宿主机中,因为已经往registry发送了,保留有些多余,需要有删除机制.
 * jenkins没有回滚机制,部署不在jenkins这端做了,需要再找方案.
