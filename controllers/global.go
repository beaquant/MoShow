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

//TokenRedisKey .
const TokenRedisKey = "token"

var (
	key        = []byte(beego.AppConfig.String("aesKey"))
	cookieName = "tk"
	r          *regexp.Regexp
)

//Token .
type Token struct {
	ID         uint64
	AcctStatus int
	ExpireTime time.Time
}

//FreqRecord 记录某个ip在某个时间段内的调用频率
type FreqRecord struct {
	ExpireTime time.Time
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
	exclude["/api/auth/*"] = struct{}{}
	exclude["/api/order/verify"] = struct{}{}

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
		if err := utils.JSONUnMarshal(val, f); err == nil && f.ExpireTime.After(time.Now()) {
			if f.Count > count {
				dto := &utils.ResultDTO{Sucess: false, Message: "您访问频率太快，请稍后再试", Code: utils.DtoStatusFrequencyError}
				ctx.Output.JSON(dto, false, false)
				return false
			}
			f.Count++
		} else {
			f.Count = 1
			f.ExpireTime = time.Now().Add(time.Minute)
		}
	} else {
		f.Count = 1
		f.ExpireTime = time.Now().Add(time.Minute)
	}

	if str, err := utils.JSONMarshalToString(f); err == nil {
		con.Do("HSET", FrequencyRedisKey, ip, str)
	}
	return true
}

//SetToken 在cookie里添加token字段
func SetToken(ctx *context.Context, tk *Token) error {
	tkStr, err := tk.Encrypt()
	if err != nil {
		return err
	}

	ctx.SetCookie(cookieName, tkStr)
	return nil
}

//GetToken 校验token失败时，直接返回错误
func GetToken(ctx *context.Context) *Token {
	ckStr := ctx.GetCookie(cookieName)

	b := &Token{}
	err := b.Decrypt(ckStr)

	if err != nil {
		beego.Error(err)

		dto := &utils.ResultDTO{Sucess: false, Message: "Token校验失败,请先登录", Code: utils.DtoStatusAuthError}
		ctx.Output.JSON(dto, false, false)
		return nil
	}

	if b.ExpireTime.Before(time.Now()) {
		dto := &utils.ResultDTO{Sucess: false, Message: "Token已过期,请重新登录", Code: utils.DtoStatusAuthError}
		ctx.Output.JSON(dto, false, false)
		return nil
	}

	if b.AcctStatus == models.AcctStatusDeleted {
		dto := &utils.ResultDTO{Sucess: false, Message: "您的账号已被注销,请联系客服", Code: utils.DtoStatusAuthError}
		ctx.Output.JSON(dto, false, false)
		return nil
	}
	return b
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
