module github.com/qeelyn/golang-starter-kit

require (
	github.com/bradfitz/gomemcache v0.0.0-20180627213657-7886924cd2b3 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/gzip v0.0.0-20180827025821-1c01330167a8
	github.com/gin-gonic/gin v1.3.0
	github.com/golang/protobuf v1.1.0
	github.com/graph-gophers/dataloader v0.0.0-20180104184831-78139374585c
	github.com/graph-gophers/graphql-go v0.0.0-20181008023426-a2b6fc39dd10
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/jinzhu/gorm v1.9.1
	github.com/lib/pq v0.0.0-20180523175426-90697d60dd84 // indirect
	github.com/opentracing/opentracing-go v1.0.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.8.0
	github.com/prometheus/client_golang v0.8.0
	github.com/prometheus/procfs v0.0.0-20180629160828-40f013a808ec // indirect
	github.com/qeelyn/gin-contrib v0.0.0-20180928072209-bf9bafd6c2c4
	github.com/qeelyn/go-common v0.0.0-20181009065420-2a6c43fbc313
	github.com/spf13/viper v1.1.0
	github.com/uber/jaeger-client-go v2.14.0+incompatible
	github.com/vmihailenco/msgpack v3.3.3+incompatible // indirect
	go.uber.org/zap v1.9.1
	golang.org/x/net v0.0.0-20180811021610-c39426892332
	google.golang.org/grpc v1.14.0
	gopkg.in/vmihailenco/msgpack.v3 v3.3.3 // indirect
)

replace (
	github.com/graph-gophers/graphql-go => github.com/qeelyn/graphql-go v0.0.0-20181012014650-03df3acf1181
	golang.org/x/net => github.com/golang/net v0.0.0-20180811021610-c39426892332
	golang.org/x/sys => github.com/golang/sys v0.0.0-20180810173357-98c5dad5d1a0
	golang.org/x/text => github.com/golang/text v0.3.0
	google.golang.org/appengine => github.com/golang/appengine v1.1.0
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20180808183934-383e8b2c3b9e
	google.golang.org/grpc => github.com/grpc/grpc-go v1.14.0
)
