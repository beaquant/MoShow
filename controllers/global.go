package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"encoding/base64"
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"github.com/astaxie/beego/context"

	"github.com/astaxie/beego"
)

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

func init() {
	r, _ = regexp.Compile("/api/.+?/")
	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)
}

//FilterUser .
func FilterUser(ctx *context.Context) {
	exclude := make(map[string]struct{})
	exclude["/api/auth/"] = struct{}{}

	ss := r.FindStringSubmatch(ctx.Request.RequestURI)
	if ss != nil && len(ss) > 0 {
		_, ok := exclude[ss[0]]
		if !ok {
			GetToken(ctx)
		}
	} else {
		GetToken(ctx)
	}
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
