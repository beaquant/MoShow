package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"errors"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/silenceper/wechat/oauth"
)

var timeFormat string

func init() {
	timeFormat = "2006-01-02T15:04:05.000Z"
}

//AuthController .
type AuthController struct {
	beego.Controller
}

//SendCode .
// @Title 登陆或者注册
// @Description 登陆或者注册
// @Param   phone     path    string  true        "接收验证码的手机号"
// @Success 200 {object} utils.ResultDTO
// @router /:phone/sendcode [post]
func (c *AuthController) SendCode() {
	dto := &utils.ResultDTO{Sucess: false}
	defer dto.JSONResult(&c.Controller)
	con := utils.RedisPool.Get()
	defer con.Close()

	ip := c.Ctx.Input.IP()
	val, err := con.Do("GET", ip)
	if err != nil {
		dto.Message = err.Error()
		beego.Error(err)
		return
	}

	if val != nil {
		t, err := time.Parse(timeFormat, val.(string))
		if err != nil {
			beego.Error(err)
			return
		}

		if t.After(time.Now().Add(time.Minute * 13)) {
			dto.Message = "验证码已发送，请检查手机短信"
			return
		}
	}

	num := c.Ctx.Input.Param(":phone")
	code := strconv.Itoa(utils.RandNumber(1000, 9999))

	if res, err := utils.SendMsgByAPIKey(num, code); err != nil {
		beego.Error("发送验证码失败:\t" + res + "\r\n" + err.Error())
		dto.Message = err.Error()
	} else {
		con.Do("SET", ip, time.Now().Add(time.Minute*15).Format(timeFormat), "EX", 60*15)
		con.Do("HMSET", "code", num, code)
		dto.Sucess = true
		dto.Message = "验证码发送成功"
	}

}

//Login .
// @Title 登陆或者注册
// @Description 登陆或者注册
// @Param   phone     path    	  string  true	 "手机号"
// @Param   code             formData     string  true	 "验证码"
// @Success 200 {object} utils.ResultDTO
// @router /:phone/login [post]
func (c *AuthController) Login() {
	dto := utils.ResultDTO{Sucess: false}
	defer dto.JSONResult(&c.Controller)
	con := utils.RedisPool.Get()
	defer con.Close()

	ip := c.Ctx.Input.IP()
	phoneNum := c.Ctx.Input.Param(":phone")
	code := c.GetString("code")

	if val, err := con.Do("GET", ip); val == nil || err != nil {
		if err != nil {
			beego.Error(err)
			panic(err)
		}
		panic(errors.New("验证码已过期,请重新获取"))
	} else {
		codeEx, err := con.Do("HMGET", "code", phoneNum)
		if err != nil {
			panic(err)
		}

		if codeEx.([]interface{})[0].(string) != code {
			panic(errors.New("验证码错误"))
		}
	}

	u := &models.User{PhoneNumber: phoneNum}
	if err := u.ReadFromPhoneNumber(); err != nil {
		beego.Error(err)
		panic(err)
	}

	tk := &Token{}
	if u.ID == 0 { //该手机号未注册，执行注册逻辑
		u.AcctStatus = models.AcctStatusNormal
		u.AcctType = models.AcctTypeTelephone
		u.CreatedAt = time.Now()

		if err := u.Add(); err == nil {
			tk.ID = u.ID
			dto.Sucess = true
			SetToken(c.Ctx, tk)
		} else {
			beego.Error(err)
			dto.Message = err.Error()
		}
	} else {
		if u.AcctStatus != models.AcctStatusDeleted {
			tk.ID = u.ID
			dto.Sucess = true
			SetToken(c.Ctx, tk)
		} else {
			dto.Message = "该账号已被注销"
		}
	}
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
