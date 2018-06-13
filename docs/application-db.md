配合数据库工作
=================
本套件的数据库服务组件采用的是GROM,我们即可以使用ORM方式,也可以使用RAWSQL方式.
实际在业务开发中,大部分使用的是ORM方式.请参考[gorm](http://gorm.io/docs/)

数据库访问对象
---------------
GORM默认支持的数据库如下:
* mysql/mariadb
* progreSQL

只列出实践过的,其他的Github DAO层应该都有.

### 配置
支持多数据库访问支持,具体的数据库配置请参考GORM
```
db:
  default:
    dsn: root:@tcp(localhost:3306)/test
    maxidleconns: 10
    maxopenconns: 100
    connmaxlifetime: 7200
  test:
    dsn: root:@tcp(localhost:3306)/test2
```

### 使用

* 使用默认DAO
```
    db := app.db
```
* 使用指定的DAO
```
    db := app.GetORMByName("test")
```
### 配合protobuf使用
默认生成的protobuf文件是不包括GORM的规则的.一般手动在pd.go同级目录加入.gorm.go文件,如以下文件目录
```
schemas
  - meta
    - calendar.gorm.go
    - calendar.pb.go
    - calendar.proto
```
calendar.gorm.go的内容
```go
package meta

func (Calendar) TableName() string {
	return "meta_date"
}
```
根据grom的需要,定义相关的处理方法

### 注意事项
* 链接池设置,需要根据实际环境设置,如并发要求,线上数据库支持
* 连接的生命周期,注意应该设置比数据库的小,否则会出现部分请求失败

[下一节 Graphql](graphql.md)