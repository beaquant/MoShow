package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"encoding/base64"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/garyburd/redigo/redis"
)

//SmsCodeRedisKey 短信验证码的redis哈希表key
const SmsCodeRedisKey = "code"

//FrequencyRedisKey IP调用接口频率限制的redis哈希表key
const FrequencyRedisKey = "frequency"

//TokenValidRedisKey 单端登陆唯一有效token的redis哈希表key
const TokenValidRedisKey = "token"

var (
	key        = []byte(beego.AppConfig.String("aesKey"))
	cookieName = "tk"
	r          *regexp.Regexp
)

//Token .
type Token struct {
	ID         uint64
	AcctStatus int
	UserType   int
	ExpireTime int64
	UUID       string
}

//FreqRecord 记录某个ip在某个时间段内的调用频率
type FreqRecord struct {
	ExpireTime int64
	Count      uint64
}

func init() {
	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)
}

//FilterUser .
func FilterUser(ctx *context.Context) {
	if !FrequencyCheck(ctx) {
		return
	}

	if !strings.HasPrefix(ctx.Request.RequestURI, "/api") {
		return
	}

	exclude := make(map[string]struct{})
	exclude["/api/auth/.+/sendcode"] = struct{}{}
	exclude["/api/auth/.+/login"] = struct{}{}
	exclude["/api/auth/wechatlogin"] = struct{}{}
	exclude["/api/order/verify"] = struct{}{}
	exclude["/api/dial/nmcallback"] = struct{}{}

	for k := range exclude {
		if ok, _ := regexp.MatchString(k, ctx.Request.RequestURI); ok {
			return
		}
	}
	GetToken(ctx)
}

//FrequencyCheck 检查调用频率
func FrequencyCheck(ctx *context.Context) bool {
	con := utils.RedisPool.Get()
	defer con.Close()

	count := uint64(100) //每分钟100次
	ip := ctx.Input.IP()
	f := &FreqRecord{}
	val, _ := redis.String(con.Do("HGET", FrequencyRedisKey, ip))
	if len(val) > 0 {
		if err := utils.JSONUnMarshal(val, f); err == nil && f.ExpireTime > time.Now().Unix() {
			if f.Count > count {
				dto := &utils.ResultDTO{Sucess: false, Message: "您访问频率太快，请稍后再试", Code: utils.DtoStatusFrequencyError}
				ctx.Output.JSON(dto, false, false)
				return false
			}
			f.Count++
		} else {
			f.Count = 1
			f.ExpireTime = time.Now().Add(time.Minute).Unix()
		}
	} else {
		f.Count = 1
		f.ExpireTime = time.Now().Add(time.Minute).Unix()
	}

	if str, err := utils.JSONMarshalToString(f); err == nil {
		con.Do("HSET", FrequencyRedisKey, ip, str)
	}
	return true
}

//SetToken 在cookie里添加token字段
func SetToken(ctx *context.Context, tk *Token) error {
	var err error

	//设置token过期时间
	tk.ExpireTime = time.Now().AddDate(0, 0, 15).Unix()
	if tk.UUID, err = utils.UUIDBase64String(); err != nil {
		beego.Error("生成UUID失败", err, ctx.Request.UserAgent())
		tk.UUID = utils.RandStringBytesMaskImprSrc(16)
	}

	con := utils.RedisPool.Get()
	defer con.Close()

	if _, err := con.Do("HSET", TokenValidRedisKey, tk.ID, tk.UUID); err != nil {
		beego.Error("redis更新token操作失败", err, ctx.Request.UserAgent())
		return err
	}

	tkStr, err := tk.Encrypt()
	if err != nil {
		return err
	}

	(&models.User{ID: tk.ID}).UpdateLoginInfo(&models.UserLoginInfo{UserAgent: ctx.Request.UserAgent(), IPAddress: ctx.Input.IP()})

	ctx.SetCookie(cookieName, tkStr)
	return nil
}

//GetToken 校验token失败时，直接返回错误
func GetToken(ctx *context.Context) *Token {
	ckStr := ctx.GetCookie(cookieName)

	tk := &Token{}
	err := tk.Decrypt(ckStr)
	dto := &utils.ResultDTO{Sucess: false, Code: utils.DtoStatusAuthError}

	if err != nil {
		beego.Error("token解密失败", err)
		dto.Message = "Token校验失败,请先登录"
		ctx.Output.JSON(dto, false, false)
		return nil
	}

	if tk.ExpireTime < time.Now().Unix() {
		beego.Error("Token已过期", tk.ExpireTime, tk)
		dto.Message = "Token已过期,请重新登录"
		ctx.Output.JSON(dto, false, false)
		return nil
	}

	if tk.AcctStatus == models.AcctStatusDeleted {
		dto.Message = "您的账号已被注销,请联系客服"
		ctx.Output.JSON(dto, false, false)
		return nil
	}

	con := utils.RedisPool.Get()
	defer con.Close()
	val, err := redis.String(con.Do("HGET", TokenValidRedisKey, tk.ID))
	if err != nil {
		beego.Error("redis获取token失败", err)
		dto.Message = "校验单端登陆失败"
		ctx.Output.JSON(dto, false, false)
		return nil
	}

	if val != tk.UUID {
		dto.Message = "已在其他设备登陆,请注意信息安全"
		ctx.Output.JSON(dto, false, false)
		return nil
	}

	return tk
}

//IsCheckMode .
func IsCheckMode(agent string) bool {
	regs, err := (&models.Config{}).GetcheckModeRegs()
	if err != nil {
		beego.Error("获取审核模式正则失败", err)
		return false
	}

	for index := range regs {
		if ok, err := regexp.MatchString(regs[index], agent); err != nil {
			beego.Error("审核模式正则匹配异常", err)
			return false
		} else if ok {
			return true
		}
	}

	return false
}

//IsCheckMode4Context .
func IsCheckMode4Context(ctx *context.Context) bool {
	if ctx.Request.Header.Get("Azwx") == "0" {
		return true
	}

	if strings.Contains(strings.ToLower(ctx.Request.UserAgent()), "ipad") {
		return true
	}

	timeZone := strings.ToLower(ctx.Request.Header.Get("Client"))
	if strings.Contains(timeZone, "us") || strings.Contains(timeZone, "america") {
		if ii, err := utils.GetIPInfo(ctx.Input.IP()); err != nil && ii.CountryCode == "US" {
			return true
		}
	}

	return IsCheckMode(ctx.Request.UserAgent())
}

//ClearToken 清除token字段
func ClearToken(ctr *beego.Controller) {
	ctr.Ctx.SetCookie(cookieName, "")
}

//Encrypt .
func (tk *Token) Encrypt() (string, error) {
	str, err := json.Marshal(tk)
	if err != nil {
		return "", err
	}

	data, err := utils.AesEncrypt(str, key)
	if err != nil {
		return "", err
	}
	result := base64.URLEncoding.EncodeToString(data)
	return result, nil
}

//Decrypt .
func (tk *Token) Decrypt(str string) error {
	if len(str) == 0 {
		return errors.New("未能读取token")
	}

	result, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return err
	}

	origData, err := utils.AesDecrypt(result, key)
	if err != nil {
		return err
	}

	err = json.Unmarshal(origData, tk)
	if err != nil {
		return err
	}
	return nil
}
