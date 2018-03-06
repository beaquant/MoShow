// @APIVersion 0.1.0
// @Title MoShow Api
// @Description api的所有返回结果均为json格式,{"Sucess":true,"Data":{},"Message":"","Code":0},所有api的结果，只有Data的类型根据API改变,其他字段均不变
// @Contact mrsong0607@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"MoShow/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&controllers.AuthController{},
			),
		), beego.NSNamespace("/ws",
			beego.NSInclude(
				&controllers.WebsocketController{},
			),
		),
		beego.NSNamespace("/config",
			beego.NSInclude(
				&controllers.ConfigController{},
			),
		),
		beego.NSNamespace("/im",
			beego.NSInclude(
				&controllers.ImController{},
			),
		),
		beego.NSNamespace("/order",
			beego.NSInclude(
				&controllers.OrderController{},
			),
		),
		beego.NSNamespace("/timeline",
			beego.NSInclude(
				&controllers.TimelineController{},
			),
		),
		beego.NSNamespace("/feedback",
			beego.NSInclude(
				&controllers.FeedbackController{},
			),
		),
		beego.NSNamespace("/blchg",
			beego.NSInclude(
				&controllers.BalanceChgController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
