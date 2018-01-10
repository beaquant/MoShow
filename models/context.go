package models

import (
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:1314melodysong@/MoShow?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	beego.Info("数据库连接初始化完成")
}
