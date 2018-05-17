package utils

import (
	"errors"
	"net/http"

	"github.com/astaxie/beego"

	"github.com/smartwalle/alipay"
)

var (
	client          *alipay.AliPay
	alipayAppID     = beego.AppConfig.String("alipayAppid")
	alipayParterID  = beego.AppConfig.String("alipayParterid")
	alipayPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApNynDKyB8zcngtqnJKikIoMM8y0K+EziCPuQC1qkdYul59tq5sOmBd/VzppbnS+CdFeG3vlS02AhpP1ALmspokXeek5sliB3bwtm/cguGpy5iKAIAkL+xgGSFNeY3waiU2c9PhIGGIrT4zKoSJYta2vEc2bMAyaCAVLcYWa8tIhgiQYXmnviN2YspdXcPaYFdEmWNyUOBFVvb6eSzArzDHodn1Kyh5X4k8qshTCQ9yvpFLzJ2ZUYkI/NPIzYEjNjCEQ3F3OYThVLs0BsuVs6HcP/p13UiorT9HATp3alvhLhfdvkyITHc4zHYB4jDMY4AJBiL3QfpJ4nebDvC4pfZQIDAQAB
-----END PUBLIC KEY-----`)
	publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApe58AH1bQgbxjKOBA6IiaKZx7NqxNTgrd1BpteWRlRzDRelqIbYiUF0sok1HCxvC2kKiul2j9GfNrNZZ+qBuDXayd1IZJLyq5xnU5zSb5+QYmIWZEh7ahMeHvemHq1rYxgej0jfM7HYp9rVCQ74aJMpwlAWuP2+43QEivtsUnEYz9rQYQVHDA69QrTs2ViSIJXk4Ag2KmgjZ0ysmr7GB2nsl7YZXHSULXmJ+m7ibDvEk8EXbZBpj9e29wDP7tYiInjvLhrEJD0h+LQ1jTJTk88DGia7wPbWDZyPB/tGptkV+x9XfWcY+bUNGaAJZIQiLw27jM4hT115PWmNqJX10kQIDAQAB
-----END PUBLIC KEY-----`)
	privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEApe58AH1bQgbxjKOBA6IiaKZx7NqxNTgrd1BpteWRlRzDRelqIbYiUF0sok1HCxvC2kKiul2j9GfNrNZZ+qBuDXayd1IZJLyq5xnU5zSb5+QYmIWZEh7ahMeHvemHq1rYxgej0jfM7HYp9rVCQ74aJMpwlAWuP2+43QEivtsUnEYz9rQYQVHDA69QrTs2ViSIJXk4Ag2KmgjZ0ysmr7GB2nsl7YZXHSULXmJ+m7ibDvEk8EXbZBpj9e29wDP7tYiInjvLhrEJD0h+LQ1jTJTk88DGia7wPbWDZyPB/tGptkV+x9XfWcY+bUNGaAJZIQiLw27jM4hT115PWmNqJX10kQIDAQABAoIBAQCNAf50uBKuKJJqrqO7f7P39MJJwornLAWcDkTXI/C0o826AqKDZXEBlDyS7FLcOMo8inYZI+xpjTD2mO28E0uSu3Tr+2OMmZwuagBIPqfixy6zpoyvHnMadSmAlQ0K7Ffc6a8ovOyzYbNFiUF3qfwzmalT0QQDuqCBhy6MvEZmic9WMO/H7uE7+HimFJ63R/UT/s9TkTUSIcvHjecTF6Rx0iFKk6ESUbVjijzGbOZC3isKluRywWWTJYQBirw+CWBnKy0egJImVpYW6GZ8z4yPEO4tm4GLn88U6ITa0+vo4jtP89By7Ewb6t93bAUIhwEtOrpUKCoznxoihVYjEZDVAoGBANnmmiaBRBOcu1r9BxlXalsD7XSW8gzW01b6fbrI4IzLpVeFyG6LWkLdWpOFV0WrWkCV+48GM1Jq681oHiun1ByuEiMg5vLVXA9V/1XpLa/pKrVOTMbUchNVBZY9vD5/PkqEfiWwsBBdcvtv5un/O35UboOLqSnvmNhixkVSh7cvAoGBAMLxs9zQK7kuDE4iUqPxezSXcgQjHf0oEaAqqz+96YKNxxGnvuJzncPjS3WvvlyOZVYPgzkHyRZ6LvnUzVay8jPwg5W8ndUDqor9/G+q9poWEPriohmjue+DeiVb1bOZVLrf3KO1TRFwBWoQYmK0cBZH8cD54L5H6I52W+2vhqA/AoGBAIzQaW3Yu5WxA6KZQa0uwJxwvVNK+MEzUwAygG3kwrg6Im+dFRnbFEmBorcSxINRaNG0Gw0ihKgOULQ9RMIRgxHFrBLngFgNaaC/gnKSbePwWpkwMI2NXOsBVsrwumXo9OhTFvJkbGMnANdcSW2Oc3QAPCrmZjujirLLojXKT8ohAoGAPz7HSaZH6SYlW9wKz6FyhVd06B60hgNP5JSzRlTIw1BX+0Rey30S/BBr1NyVd9XCzq7ttbzu4ln1j5wYmj4oEe2/4v50fj1YQQuxsFDY/JiYHa0VRhg2JJyVLjWjGUdvk8k4/eu9+yBKwWRbZwZ/LttcdW0cGt+ddUq0/GHr3WUCgYAUlojAOI7zHEvvjNS8/DIHACsiv06Z+Hpb8B2NJy3aolIZklft+pj4CODgM1IyxRLsyXLV0oH8E/naoLN1bUwjJ84e6gmIIuUipyvyeCwWjAxlyCICRQyQ4zvU1WDJ1Sqikky0PV7S58Ot/dziI0bjXnKOY6dPOebT03NZo7sc3A==
-----END RSA PRIVATE KEY-----`)
)

func init() {
	client = alipay.New(alipayAppID, alipayParterID, publicKey, privateKey, true)
	client.AliPayPublicKey = alipayPublicKey
}

//CreatePayment .
func CreatePayment(title, orderID, NotifyURL, Amount string) (string, error) {
	var p = alipay.AliPayTradeWapPay{}
	p.NotifyURL = NotifyURL
	p.ReturnURL = "http://xinmeiwl.com/"
	p.Subject = title
	p.OutTradeNo = orderID
	p.TotalAmount = Amount
	p.ProductCode = "QUICK_WAP_WAY"

	var url, err = client.TradeWapPay(p)
	if url == nil {
		return "", err
	}
	return url.String(), err
}

//ConfirmPayment .
func ConfirmPayment(req *http.Request) (*alipay.TradeNotification, error) {
	notify, err := client.GetTradeNotification(req)
	if err != nil {
		return notify, err
	}

	if notify.AppId != alipayAppID {
		return notify, errors.New("异常通知:appId检验失败")
	}

	if notify.SellerId != alipayParterID {
		return notify, errors.New("异常通知:seller_id检验失败")
	}

	if notify.TradeStatus != alipay.K_ALI_PAY_TRADE_STATUS_TRADE_SUCCESS && notify.TradeStatus != alipay.K_ALI_PAY_TRADE_STATUS_TRADE_FINISHED {
		return nil, nil
	}

	return notify, nil
}
