package utils

import (
	"errors"
	"fmt"

	"github.com/MrSong0607/aliyun-communicate"
	"github.com/astaxie/beego"
)

var (
	smsAccessKeyID, smsAccessKeySecret string
	gatewayURL                         = "http://dysmsapi.aliyuncs.com/"
	signName                           = "蜜秀MX"
	templateCode                       = "SMS_130924577"
	templateParam                      = "{\"code\":\"%s\"}"
)

func init() {
	smsAccessKeyID = beego.AppConfig.String("smsAccessKeyId")
	smsAccessKeySecret = beego.AppConfig.String("smsAccessKeySecret")
}

//SendMsgByAPIKey .
func SendMsgByAPIKey(mobile, content string) (string, error) {
	smsClient := aliyunsmsclient.New(gatewayURL)
	result, err := smsClient.Execute(smsAccessKeyID, smsAccessKeySecret, mobile, signName, templateCode, fmt.Sprintf(templateParam, content))
	if err != nil {
		return "", err
	}

	rs, err := JSONMarshalToString(result)
	if result.Code != "OK" {
		return "", errors.New(rs)
	}

	return rs, err
}
