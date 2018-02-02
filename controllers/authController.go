package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"github.com/silenceper/wechat/oauth"
)

var timeFormat = "2006-01-02T15:04:05.000Z"
var adminPhone = beego.AppConfig.String("adminPhoneNum")
var adminCode = beego.AppConfig.String("adminCode")

//AuthController 短信登陆，微信登陆，发送验证码，退出登陆等
type AuthController struct {
	beego.Controller
}

type codeInfo struct {
	Code string
	Time time.Time
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
	val, _ := redis.String(con.Do("GET", ip))

	num := c.Ctx.Input.Param(":phone")
	if len(val) > 0 {
		if num == val {
			dto.Message = "验证码已发送，请检查手机短信"
		} else {
			dto.Message = "验证码请求太频繁，请稍等"
		}
		return
	}

	code := strconv.Itoa(utils.RandNumber(1000, 9999))

	if res, err := utils.SendMsgByAPIKey(num, code); err != nil {
		beego.Error("发送验证码失败:\t" + res + "\r\n" + err.Error())
		dto.Message = err.Error()
	} else {
		cs, _ := utils.JSONMarshalToString(&codeInfo{Code: code, Time: time.Now().Add(time.Minute * 15)})

		con.Do("SET", ip, num, "EX", 60*2)
		con.Do("HMSET", "code", num, cs)
		dto.Sucess = true
		dto.Message = "验证码发送成功"
	}
}

//Login .
// @Title 登陆或者注册
// @Description 登陆或者注册
// @Param   phone     path    	  string  true	 "手机号"
// @Param   code	  formData     string  true	 "验证码"
// @Success 200 {object} utils.ResultDTO
// @router /:phone/login [post]
func (c *AuthController) Login() {
	dto := utils.ResultDTO{Sucess: false}
	defer dto.JSONResult(&c.Controller)
	con := utils.RedisPool.Get()
	defer con.Close()

	phoneNum := c.Ctx.Input.Param(":phone")
	code := c.GetString("code")

	if phoneNum != adminPhone && code != adminCode {
		codeEx, err := redis.Strings(con.Do("HMGET", "code", phoneNum))
		if err != nil {
			beego.Error(err)
			dto.Message = err.Error()
			return
		}

		ci := &codeInfo{}
		utils.JSONUnMarshal(codeEx[0], ci)

		if ci.Time.Before(time.Now()) {
			dto.Message = "验证码已过期,请重新获取"
			return
		}

		if ci.Code != code {
			dto.Message = "验证码错误"
			return
		}
	}

	u := &models.User{PhoneNumber: phoneNum}
	if err := u.ReadFromPhoneNumber(); err != nil {
		beego.Error(err)
		dto.Message = err.Error()
		return
	}

	tk := &Token{ExpireTime: time.Now().AddDate(0, 0, 15)}
	if u.ID == 0 { //该手机号未注册，执行注册逻辑
		u.AcctStatus = models.AcctStatusNormal
		u.AcctType = models.AcctTypeTelephone
		u.CreatedAt = time.Now().Unix()

		if err := u.Add(); err == nil {
			tk.ID = u.ID
			dto.Message = "注册成功"
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
			dto.Message = "登陆成功"
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
	dto := utils.ResultDTO{Sucess: false}
	defer dto.JSONResult(&c.Controller)

	AccessToken := c.GetString("AccessToken")
	OpenID := c.GetString("OpenID")

	o := oauth.NewOauth(nil)
	info, err := o.GetUserInfo(AccessToken, OpenID)
	if err != nil {
		dto.Message = err.Error()
		beego.Error(err)
		return
	}

	u := &models.User{WeChatID: info.OpenID}
	err = u.ReadFromWechatID()
	if err != nil {
		dto.Message = err.Error()
		beego.Error(err)
		return
	}

	tk := &Token{}
	if u.ID == 0 { //执行微信注册
		u.AcctType = models.AcctTypeWechat
		u.AcctStatus = models.AcctStatusNormal
		u.CreatedAt = time.Now().Unix()
		u.Add()

		if err := u.Add(); err == nil {
			tk.ID = u.ID
			dto.Message = "注册成功"
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
			dto.Message = "登陆成功"
			SetToken(c.Ctx, tk)
		} else {
			dto.Message = "该账号已被注销"
		}
	}
}

//Logout .
// @Title 注销登录
// @Description 注销登录
// @Success 200 {object} utils.ResultDTO
// @router /logout [get]
func (c *AuthController) Logout() {
	ClearToken(&c.Controller)
	dto := utils.ResultDTO{Sucess: false, Message: "退出登陆成功"}
	dto.JSONResult(&c.Controller)
}
