package main

import (
	_ "github.com/go-sql-driver/mysql"

	_ "MoShow/routers"

	"github.com/astaxie/beego"
)

func init() {
	// beego.BConfig.RunMode = "dev"
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
