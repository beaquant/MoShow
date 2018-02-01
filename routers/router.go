// @APIVersion 0.1.0
// @Title MoShow Api
// @Description beego has a very cool tools to autogenerate documents for your API
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
	)
	beego.AddNamespace(ns)
}
