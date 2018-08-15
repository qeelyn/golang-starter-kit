module github.com/qeelyn/golang-starter-kit

require (
	github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20180627213657-7886924cd2b3 // indirect
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/coreos/etcd v3.3.8+incompatible // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7 // indirect
	github.com/gin-gonic/gin v0.0.0-20170702092826-d459835d2b07
	github.com/gogo/protobuf v1.0.0 // indirect
	github.com/golang/protobuf v1.1.0
	github.com/graph-gophers/dataloader v0.0.0-20180104184831-78139374585c
	github.com/graph-gophers/graphql-go v0.0.0-20180604122119-0b810f691a45
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/hcl v0.0.0-20180404174102-ef8a98b0bbce // indirect
	github.com/jinzhu/gorm v1.9.1
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/lib/pq v0.0.0-20180523175426-90697d60dd84 // indirect
	github.com/magiconair/properties v1.8.0 // indirect
	github.com/mattn/go-isatty v0.0.3 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/mapstructure v0.0.0-20180511142126-bb74f1db0675 // indirect
	github.com/opentracing/opentracing-go v1.0.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/pkg/errors v0.8.0
	github.com/prometheus/client_golang v0.8.0
	github.com/prometheus/client_model v0.0.0-20171117100541-99fa1f4be8e5 // indirect
	github.com/prometheus/common v0.0.0-20180518154759-7600349dcfe1 // indirect
	github.com/prometheus/procfs v0.0.0-20180629160828-40f013a808ec // indirect
	github.com/qeelyn/gin-contrib v0.0.0-20180713070920-78e93a6803dc
	github.com/qeelyn/go-common v0.0.0-20180717161004-65a10cfbdf71
	github.com/spf13/afero v1.1.1 // indirect
	github.com/spf13/cast v1.2.0 // indirect
	github.com/spf13/jwalterweatherman v0.0.0-20180109140146-7c0cea34c8ec // indirect
	github.com/spf13/pflag v1.0.1 // indirect
	github.com/spf13/viper v1.0.2
	github.com/uber/jaeger-client-go v2.14.0+incompatible
	github.com/uber/jaeger-lib v1.5.0 // indirect
	github.com/ugorji/go v1.1.1 // indirect
	github.com/vmihailenco/msgpack v3.3.3+incompatible // indirect
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.8.0
	golang.org/x/net v0.0.0-20180629035331-4cb1c02c05b0
	golang.org/x/sys v0.0.0-20180627142611-7138fd3d9dc8 // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/appengine v1.1.0 // indirect
	google.golang.org/genproto v0.0.0-20180627194029-ff3583edef7d // indirect
	google.golang.org/grpc v1.13.0
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0-20170531160350-a96e63847dc3 // indirect
	gopkg.in/vmihailenco/msgpack.v3 v3.3.3 // indirect
	gopkg.in/yaml.v2 v2.2.1 // indirect
)

replace (
	github.com/bradfitz/gomemcache => gitee.com/githubmirror/gomemcache v0.0.0-20180710155616-bc664df96737
	github.com/coreos/etcd => gitee.com/githubmirror/etcd v3.3.9+incompatible
	github.com/dgrijalva/jwt-go => gitee.com/githubmirror/jwt-go v3.2.0+incompatible
	github.com/graph-gophers/graphql-go => github.com/qeelyn/graphql-go v0.0.0-20180604122119-0b810f691a45
	github.com/uber/jaeger-client-go => gitee.com/githubmirror/jaeger-client-go v2.14.0+incompatible
	github.com/vmihailenco/msgpack => gitee.com/githubmirror/msgpack v3.3.3+incompatible
	golang.org/x/net => github.com/golang/net v0.0.0-20180811021610-c39426892332
	golang.org/x/sys => github.com/golang/sys v0.0.0-20180810173357-98c5dad5d1a0
	golang.org/x/text => github.com/golang/text v0.3.0
	google.golang.org/appengine => gitee.com/githubmirror/appengine v1.1.0
	google.golang.org/genproto => gitee.com/githubmirror/go-genproto v0.0.0-20180808183934-383e8b2c3b9e
	google.golang.org/grpc => gitee.com/githubmirror/grpc-go v1.14.0
)
