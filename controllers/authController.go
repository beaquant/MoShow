package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"strconv"
	"time"

	netease "github.com/MrSong0607/netease-im"
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
// @Title 发送验证码
// @Description 发送验证码
// @Param   phone     path    string  true        "接收验证码的手机号"
// @Success 200 {object} utils.ResultDTO
// @router /:phone/sendcode [post]
func (c *AuthController) SendCode() {
	dto := &utils.ResultDTO{Sucess: false}
	defer dto.JSONResult(&c.Controller)
	con := utils.RedisPool.Get()
	defer con.Close()

	num := c.Ctx.Input.Param(":phone")
	codeEx, err := redis.String(con.Do("HGET", SmsCodeRedisKey, num))
	if err != nil {
		beego.Error(err)
		dto.Message = err.Error()
		return
	}

	ci := &codeInfo{}
	utils.JSONUnMarshal(codeEx, ci)

	if ci != nil && ci.Time.After(time.Now().Add(time.Minute*13)) {
		dto.Message = "验证码请求太频繁，请稍等"
		return
	}

	code := strconv.Itoa(utils.RandNumber(1000, 9999))

	if res, err := utils.SendMsgByAPIKey(num, code); err != nil {
		beego.Error("发送验证码失败:\t" + res + "\r\n" + err.Error())
		dto.Message = err.Error()
	} else {
		cs, _ := utils.JSONMarshalToString(&codeInfo{Code: code, Time: time.Now().Add(time.Minute * 15)})

		con.Do("HSET", SmsCodeRedisKey, num, cs)
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
		codeEx, err := redis.String(con.Do("HGET", SmsCodeRedisKey, phoneNum))
		if err != nil {
			beego.Error(err)
			dto.Message = err.Error()
			return
		}

		ci := &codeInfo{}
		utils.JSONUnMarshal(codeEx, ci)

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
		trans := models.TransactionGen()

		u.AcctType = models.AcctTypeTelephone
		u.AcctStatus = models.AcctStatusNormal
		u.CreatedAt = time.Now().Unix()

		if err := u.Add(trans); err != nil {
			beego.Error(err, c.Ctx.Request.UserAgent())
			dto.Message = err.Error()
			models.TransactionRollback(trans)
			return
		}

		imUser := &netease.ImUser{ID: strconv.FormatUint(u.ID, 10)}
		imtk, err := utils.ImCreateUser(imUser)
		if err != nil {
			beego.Error("创建IMUser失败", err, c.Ctx.Request.UserAgent())
			dto.Message = "创建IMUser失败\t" + err.Error()
			models.TransactionRollback(trans)
			return
		}

		up := models.UserProfile{ID: u.ID}
		up.ImToken = imtk.Token
		up.Birthday = time.Date(1993, 1, 1, 0, 0, 0, 0, time.Local).Unix()
		up.Following = "{}"
		up.Followers = "{}"
		up.CoverPic = "{}"
		if err := up.Add(trans); err != nil {
			beego.Error(err, c.Ctx.Request.UserAgent())
			dto.Message = err.Error()
			models.TransactionRollback(trans)
			return
		}

		models.TransactionCommit(trans)
		tk.ID = u.ID
		dto.Message = "注册成功"
		dto.Data = &UserPorfileInfo{UserProfile: up, ImTk: imtk.Token}
		dto.Sucess = true
		SetToken(c.Ctx, tk)
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
		trans := models.TransactionGen()

		u.AcctType = models.AcctTypeWechat
		u.AcctStatus = models.AcctStatusNormal
		u.CreatedAt = time.Now().Unix()

		if err := u.Add(trans); err != nil {
			beego.Error(err, c.Ctx.Request.UserAgent())
			dto.Message = err.Error()
			models.TransactionRollback(trans)
			return
		}

		imUser := &netease.ImUser{ID: strconv.FormatUint(u.ID, 10)}
		imtk, err := utils.ImCreateUser(imUser)
		if err != nil {
			beego.Error("创建IMUser失败", err, c.Ctx.Request.UserAgent())
			dto.Message = "创建IMUser失败\t" + err.Error()
			models.TransactionRollback(trans)
			return
		}

		up := models.UserProfile{ID: u.ID}
		up.ImToken = imtk.Token
		up.Birthday = time.Date(1993, 1, 1, 0, 0, 0, 0, time.Local).Unix()
		up.Following = "{}"
		up.Followers = "{}"
		up.CoverPic = "{}"
		if err := up.Add(trans); err != nil {
			beego.Error(err, c.Ctx.Request.UserAgent())
			dto.Message = err.Error()
			models.TransactionRollback(trans)
			return
		}

		models.TransactionCommit(trans)
		tk.ID = u.ID
		dto.Message = "注册成功"
		dto.Data = &UserPorfileInfo{UserProfile: up, ImTk: imtk.Token}
		dto.Sucess = true
		SetToken(c.Ctx, tk)
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
