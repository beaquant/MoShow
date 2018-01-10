package controllers

import (
	"github.com/astaxie/beego"
	"github.com/silenceper/wechat/oauth"
)

//AuthController .
type AuthController struct {
	beego.Controller
}

//Login .
// @Title 登陆或者注册
// @Description 登陆或者注册
// @Param   key     path    string  true        "The email for login"
// @Success 200 {object} utils.ResultDTO
// @router /login [post]
func (c *AuthController) Login() {
}

//WechatLogin .
// @Title 登陆或者注册
// @Description 登陆或者注册
// @Param   AccessToken     formData     string  true        "The email for login"
// @Param   OpenID          formData     string  true        "The email for login"
// @Success 200 {object} utils.ResultDTO
// @router /wechatlogin [post]
func (c *AuthController) WechatLogin() {
	AccessToken := c.GetString("AccessToken")
	OpenID := c.GetString("OpenID")

	o := oauth.NewOauth(nil)
	o.GetUserInfo(AccessToken, OpenID)
}

//Logout .
// @Title 注销登录
// @Description 注销登录
// @Success 200 {object} utils.ResultDTO
// @router /logout [get]
func (c *AuthController) Logout() {

}
