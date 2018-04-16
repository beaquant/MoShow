package utils

import (
	"errors"
	"fmt"

	"github.com/KenmyZhang/aliyun-communicate/app"
	"github.com/astaxie/beego"
)

var (
	smsAccessKeyID, smsAccessKeySecret string
	gatewayURL                         = "http://dysmsapi.aliyuncs.com/"
	signName                           = "美拉拉"
	templateCode                       = "SMS_107930001"
	templateParam                      = "{\"code\":\"%s\"}"
)

func init() {
	smsAccessKeyID = beego.AppConfig.String("smsAccessKeyId")
	smsAccessKeySecret = beego.AppConfig.String("smsAccessKeySecret")
}

//SendMsgByAPIKey .
func SendMsgByAPIKey(mobile, content string) (string, error) {
	smsClient := app.NewSmsClient(gatewayURL)
	result, err := smsClient.Execute(smsAccessKeyID, smsAccessKeySecret, mobile, signName, templateCode, fmt.Sprintf(templateParam, content))
	if err != nil {
		return "", err
	}

	rs, err := JSONMarshalToString(result)
	if val, ok := result["Code"].(string); ok && val != "OK" {
		return "", errors.New(rs)
	}

	return rs, err
}
