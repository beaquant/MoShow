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

	beego.GlobalControllerRouter["MoShow/controllers:BalanceChgController"] = append(beego.GlobalControllerRouter["MoShow/controllers:BalanceChgController"],
		beego.ControllerComments{
			Method: "GetChgDetail",
			Router: `/:chgid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:BalanceChgController"] = append(beego.GlobalControllerRouter["MoShow/controllers:BalanceChgController"],
		beego.ControllerComments{
			Method: "GetIncomeList",
			Router: `/incomes`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:BalanceChgController"] = append(beego.GlobalControllerRouter["MoShow/controllers:BalanceChgController"],
		beego.ControllerComments{
			Method: "GetPaymentList",
			Router: `/payments`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:ConfigController"] = append(beego.GlobalControllerRouter["MoShow/controllers:ConfigController"],
		beego.ControllerComments{
			Method: "GetCosSign",
			Router: `/cossign`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:ConfigController"] = append(beego.GlobalControllerRouter["MoShow/controllers:ConfigController"],
		beego.ControllerComments{
			Method: "GetGiftList",
			Router: `/gifts`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:ConfigController"] = append(beego.GlobalControllerRouter["MoShow/controllers:ConfigController"],
		beego.ControllerComments{
			Method: "GetProductList",
			Router: `/products`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:DialController"] = append(beego.GlobalControllerRouter["MoShow/controllers:DialController"],
		beego.ControllerComments{
			Method: "Del",
			Router: `/:dialid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:DialController"] = append(beego.GlobalControllerRouter["MoShow/controllers:DialController"],
		beego.ControllerComments{
			Method: "DialList",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:FeedbackController"] = append(beego.GlobalControllerRouter["MoShow/controllers:FeedbackController"],
		beego.ControllerComments{
			Method: "Suggestion",
			Router: `/suggestion`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:ImController"] = append(beego.GlobalControllerRouter["MoShow/controllers:ImController"],
		beego.ControllerComments{
			Method: "CreateImUser",
			Router: `/create`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:ImController"] = append(beego.GlobalControllerRouter["MoShow/controllers:ImController"],
		beego.ControllerComments{
			Method: "RefreshToken",
			Router: `/refreshtoken`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:OrderController"] = append(beego.GlobalControllerRouter["MoShow/controllers:OrderController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/:orderid/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:OrderController"] = append(beego.GlobalControllerRouter["MoShow/controllers:OrderController"],
		beego.ControllerComments{
			Method: "AlipayConfirm",
			Router: `/verify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:OrderController"] = append(beego.GlobalControllerRouter["MoShow/controllers:OrderController"],
		beego.ControllerComments{
			Method: "CreateWebPay",
			Router: `/webpay`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:TimelineController"] = append(beego.GlobalControllerRouter["MoShow/controllers:TimelineController"],
		beego.ControllerComments{
			Method: "NewCommers",
			Router: `/users`,
			AllowHTTPMethods: []string{"get"},
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
			Method: "Follow",
			Router: `/:userid/follow`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "GiftHistory",
			Router: `/:userid/gifthis`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "Report",
			Router: `/:userid/report`,
			AllowHTTPMethods: []string{"post"},
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
			Method: "UnFollow",
			Router: `/:userid/unfollow`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "AnchorApply",
			Router: `/acapply`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "ReduceAmount",
			Router: `/cutamount`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetFollowedLst",
			Router: `/fanslist`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "GuestList",
			Router: `/guests`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "InviteList",
			Router: `/ivtlist`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "SetBusyStatus",
			Router: `/setbusy`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetFollowingLst",
			Router: `/sublist`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/update`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "Withdraw",
			Router: `/withdraw`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MoShow/controllers:UserController"] = append(beego.GlobalControllerRouter["MoShow/controllers:UserController"],
		beego.ControllerComments{
			Method: "WithdrawHis",
			Router: `/withdrawhis`,
			AllowHTTPMethods: []string{"get"},
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
			Method: "Reject",
			Router: `/:channelid/reject`,
			AllowHTTPMethods: []string{"post"},
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
