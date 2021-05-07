package model

import (
	"litrocket/utils/handlelog"
	"time"

	. "litrocket/common"

	_ "github.com/go-sql-driver/mysql" //mysql driver
	"github.com/jinzhu/gorm"           //gorm
)

var Db *gorm.DB
var err error

//InitDb
func InitDb() {
	// Connect to db.
	Db, err = gorm.Open("mysql", DbUser+":"+DbPass+"@("+DbHost+":"+DbPort+")/"+DbName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		handlelog.Handlelog("FATAL", "Model:db InitDb Failed."+err.Error())
	}

	//设置自动迁移这些表,在新的环境部署时，会自动创建表和这些表的字段，属性等。
	//但是不在新环境的话，修改结构体的字段更改不到数据库(自己使用语句改，或者删除表再重新初始化)
	Db.AutoMigrate(&Article{}, &Friend{}, &Group{}, &GroupInfo{}, &Message{}, &User{}, &File{})

	//设置禁用复数,不设置的话自动将表名user改为users, 设置了但是没用，表名还是自动搞成复数了,可能是新版本已经不让修改,而且现在加上这个会有错误，操作表时报错表找不到
	//Db.SingularTable(true)

	//设置连接池中最大闲置连接数
	Db.DB().SetMaxIdleConns(10)

	//设置数据库最大连接数量
	Db.DB().SetMaxOpenConns(100)

	//设置连接最大可复用时间
	Db.DB().SetConnMaxLifetime(10 * time.Second)

	handlelog.Handlelog("INFO", "Init Database ok")
}
