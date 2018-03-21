package utils

import (
	"errors"
	"net/http"

	"github.com/smartwalle/alipay"
)

var client *alipay.AliPay
var alipayAppID = "2017091508742831"
var alipayParterId = "2088821012806925"

func init() {
	client = alipay.New(alipayAppID, alipayParterId, publicKey, privateKey, true)
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

	if notify.SellerId != alipayParterId {
		return notify, errors.New("异常通知:seller_id检验失败")
	}

	if notify.TradeStatus != alipay.K_ALI_PAY_TRADE_STATUS_TRADE_SUCCESS && notify.TradeStatus != alipay.K_ALI_PAY_TRADE_STATUS_TRADE_FINISHED {
		return nil, nil
	}

	return notify, nil
}

// RSA2(SHA256)
var (
	alipayPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAr8nWcRBXLF8OKDV0w4EtQSmoGfQj5
h5w3RsMF8SsgXlUHmvU6Vj63Snebnmd1Fo1oflY6+nxAidwh/5P4G0aAOyS+ATUb+AqP11FXR
0f1ZJGXISA2CpJHUuN0O7hrZU33XUHaIvrYby8jDMpa9r8fnc002ZUX8elys9x+OtCmTv+ppT
xzQSf45gFMv3fvFmwxz9Sm7rSq4CgRrbnW6JzckiOYpSqClvfr+eR3c7g7A1NDpYeW0dcfbz+
Fmb/67Qs3PELJNZueW5moc8kjjK20dOYkuVcf9TAp23klHwfmiATIIDjOnYWIRlCYOPJCdrlR
7xb2GYj3xGjGHeMkJrQ/wIDAQAB
-----END PUBLIC KEY-----`)

	publicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6O8O5ei6BjQyi0+FmyfFdP/vvX49+
EpJS+/wGQVvGAKyHhDDjSFG0bnAfrNRg8ofLrsnWeu/TiZn785MtEJS5h+SvlY/y4kvbL4MYf
UL57bT6YpeivcKPDJsdAKdjb8cmwQpC38VPemJEdTYm0L/L9MmeDU+KBAgM7tIs3cszFrgKU1
Tm5xZXeTurwu/tz4ohVUAgeatwoCdNe03/2coJqcDLIdiDqC9UqnkMTLtdOSNTRXEO3gvImlv
q7MxrMu51qMWNJ3ONiWvgAzXZPuYQ7ELq5NBRhSAUejX3SWEERLsol3F3aheKoZyLBEXkMHDl
dBJ+M6eY2heAy6U1uSkmwIDAQAB
-----END PUBLIC KEY-----`)

	privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA6O8O5ei6BjQyi0+FmyfFdP/vvX49+EpJS+/wGQVvGAKyHhDDjSFG
0bnAfrNRg8ofLrsnWeu/TiZn785MtEJS5h+SvlY/y4kvbL4MYfUL57bT6YpeivcKPDJs
dAKdjb8cmwQpC38VPemJEdTYm0L/L9MmeDU+KBAgM7tIs3cszFrgKU1Tm5xZXeTurwu/
tz4ohVUAgeatwoCdNe03/2coJqcDLIdiDqC9UqnkMTLtdOSNTRXEO3gvImlvq7MxrMu5
1qMWNJ3ONiWvgAzXZPuYQ7ELq5NBRhSAUejX3SWEERLsol3F3aheKoZyLBEXkMHDldBJ
+M6eY2heAy6U1uSkmwIDAQABAoIBADLs0N0C3Giu0L4UTKl3MHw72DYde37k/tFTS7Ks
tXUSYc6g65/XPpxrd+I7Yf7mGZsl35yRJen4C5EqESr3tRgKnxJt2NSu0Wd8xUhFQq0O
E5ZjYfgPunUesQdL6xYqSh658h9JUWOPwx/a4OBQ7WmPSEHPoh8wJ1on/+T8kRa7/JfI
iI1z4aVrMLJjDVw+fzT2btb3Cff4cJCctex+cq8a5Xvh2TyM2tyJU9X1uXW7q1emNrLu
uz7N0YhHHWO/3hkQfkYskdN5+8elqeYI86kSTrMPym0jiDtFAbItxUAfWlUE+kMnCl1k
2nYUPi+zsYSAZO8D5GFU0ZMSkV8bQQECgYEA/pSgB7j7/OvNsU89y/t7n6xYo1ke3d7w
xoTra7I80u2wHqDZsMEwK2uRJQVKKMECp9V9RwNYT/gWDuTAMFFU2vDQihqB2GkAyCXa
3uE/3yTtH9qeZmEeOeS1il/1bDVJn7YsPN9D2srd4no5O5n+hroxRecdYu9gWu5Tb3ON
EYsCgYEA6juJIJyCExiGnwsebhmZvGLRs2WRDPDP3WQDQOqajZVYSzI2P7I+jXJio3RZ
suqttIYBJgHauQuTQ8uC/AR2+wonUhfqEz7alv4eN8UBKQDthRdzisL6wM1NqES8wNj8
//HHxk+zUNZz9o6PmZIZ+ogcob0zrEezDjuElv4z+zECgYEA4wsJ7dk8YsSqHYfeRR1z
k2PRaV0B+j3p3iKNEu9S74qrl6U8gDbLDu5P9ARTryTziVsM71g8WpWWlpHMFUtzsg8y
7PfW9XowCFA6cqvQmuID2HTQ792NZ3Rhs5cA+hBMKPP/YAp+KZLjcCgxAsbECMPlTcJg
out5s575KlyTYyECgYEAsdZt8KKjZ5gxbcNlYTZysMNeb5RnoqmbSH3MspbsrR58oOsI
oSfVslLsbSnDiMIBDJTJfm/d/qy5LLnxQyKoq0U0QXICuIX6NLXPf4xFqzoXG/uIMAyF
kajOkzlNDiYxQKnzga+1d2S7OrFWecShkMOS6YHbH6x4WA/8RR/Pm6ECgYEA1LuXtzyg
iLgX86KSRibepL30KDqMw9Ex5Do+7JYf+k3R9em8f154GxTR7KRfGYh6VTXKuAsR/Bjy
/zrF6vFU4kOvRQWT+3hulX1sNtM/f+Yw2VjH8BgtYYewUed6tWKMY8Wh7B9NmiV0rTmt
ahpURD9Xkef/j+/OvO8AKv89rYw=
-----END RSA PRIVATE KEY-----`)
)
