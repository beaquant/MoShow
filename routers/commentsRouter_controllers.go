package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["MoShow/controllers:AuthController"] = append(beego.GlobalControllerRouter["MoShow/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/:phone/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:AuthController"] = append(beego.GlobalControllerRouter["MoShow/controllers:AuthController"],
		beego.ControllerComments{
			Method: "SendCode",
			Router: `/:phone/sendcode`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:AuthController"] = append(beego.GlobalControllerRouter["MoShow/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:AuthController"] = append(beego.GlobalControllerRouter["MoShow/controllers:AuthController"],
		beego.ControllerComments{
			Method: "WechatLogin",
			Router: `/wechatlogin`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "Read",
			Router: `/:userid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "Del",
			Router: `/:userid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "SendGift",
			Router: `/:userid/sendgift`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/create`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/update`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:WebsocketController"] = append(beego.GlobalControllerRouter["MoShow/controllers:WebsocketController"],
		beego.ControllerComments{
			Method: "Join",
			Router: `/:channelid/join`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:WebsocketController"] = append(beego.GlobalControllerRouter["MoShow/controllers:WebsocketController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/:parterid/create`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
