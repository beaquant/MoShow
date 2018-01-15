package models

import (
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	//.
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", beego.AppConfig.String("db"))
	if err != nil {
		panic(err)
	}

	beego.Info("数据库连接初始化完成")
}
