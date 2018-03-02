package models

import (
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", beego.AppConfig.String("db"))
	if err != nil {
		panic(err)
	}

	db = db.Debug()

	beego.Info("数据库连接初始化完成")
}

//TransactionGen .
func TransactionGen() *gorm.DB {
	return db.Begin()
}

//TransactionCommit .
func TransactionCommit(trans *gorm.DB) {
	trans.Commit()
}

//TransactionRollback .
func TransactionRollback(trans *gorm.DB) {
	trans.Rollback()
}
