package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "Read",
			Router: `/:userid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/:userid`,
			AllowHTTPMethods: []string{"post"},
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
			Method: "Create",
			Router: `/create`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
