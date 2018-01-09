package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var (
	key        []byte
	cookieName string
)

func init() {
	cookieName = "tk"
	keyStr := "0123456789melody0123456789melody"
	key = []byte(keyStr)
}

//Token .
type Token struct {
	ID         uint64
	Name       string
	UserType   int
	ExpireTime time.Time
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

	if err != nil || b.ExpireTime.Before(time.Now()) || b.UserType != 0 {
		beego.Error(err)

		dto := &ResultDTO{Sucess: false, Message: "Token校验失败,请先登录", Code: StatusAuthError}
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

	data, err := AesEncrypt(str, key)
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

	origData, err := AesDecrypt(result, key)
	if err != nil {
		return err
	}

	err = json.Unmarshal(origData, tk)
	if err != nil {
		return err
	}
	return nil
}

//PKCS7Padding .
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//PKCS7UnPadding .
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AesEncrypt .
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//AesDecrypt .
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
