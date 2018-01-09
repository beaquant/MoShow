package controllers

import (
	"MoShow/utils"

	"github.com/astaxie/beego/context"

	"github.com/astaxie/beego"
)

func init() {
	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)
}

//FilterUser .
func FilterUser(ctx *context.Context) {
	if ctx.Request.RequestURI != "/user/login" {
		utils.GetToken(ctx)
	}
}
