package test

import (
	"os"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"MoShow/controllers"
	"MoShow/models"
	"MoShow/utils"
	"testing"
)

func TestCheckModePattern(t *testing.T) {
	os.Setenv("GOCACHE", "off")
	t.Log(controllers.IsCheckMode("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36 QIHU 360SE/Nutch-1.13"))
}

func TestActive(t *testing.T) {
	os.Setenv("GOCACHE", "off")
	go func() { t.Log(controllers.SendActivity(169298)) }()
	go func() { t.Log(controllers.SendActivity(169293)) }()
	go func() { t.Log(controllers.SendActivity(169143)) }()
	go func() { t.Log(controllers.SendActivity(170034)) }()
	time.Sleep(150 * time.Second)
}

func TestToken(t *testing.T) {
	tk := &controllers.Token{}
	t.Log(tk.Decrypt("iA8aheaDexLkxoijX4Yvz0vq_C52QcSmVCc6BlUpy8WpcDDZdG8s_XpO4pvAyJcIsDrxhjn7mWNJRrQJTGc0CGCMdFsWW9OiDQZAVGL7x40="))
	t.Log(utils.JSONMarshalToString(tk))
}

func TestGenWebpayURL(t *testing.T) {
	const appName = "MoShow"
	uid := uint64(169298)
	prod := &models.Product{}

	if err := utils.JSONUnMarshal(`{"name": "蜜豆充值", "extra": 0, "price": 9, "coin_count": 900, "product_id": 2}`, prod); err != nil {
		t.Error(err)
		return
	}

	u := models.User{ID: uid}
	if err := u.Read(); err != nil {
		beego.Error(err)
	} else if u.InvitedBy != 0 {
		prod.Extra++
		prod.Extra /= 2 //邀请用户充值奖励减半
	}

	prodInfo, err := utils.JSONMarshalToString(prod)

	if err != nil {
		t.Error("解析产品信息出错")
		return
	}

	trans := models.TransactionGen()
	o := &models.Order{Amount: prod.Price, CoinCount: prod.CoinCount + prod.Extra, UserID: uid, PayType: models.PayTypeAlipay, CreateAt: time.Now().Unix(), ProductInfo: prodInfo, PayInfo: "{}"}
	if err := o.Add(trans); err != nil {
		t.Error("添加订单失败\t" + err.Error())
		models.TransactionRollback(trans)
		return
	}

	//用APPName加上订单ID拼接成唯一交易号(支付宝规定每个收款账号下面的交易号必须唯一)
	uri, err := utils.CreatePayment(prod.ProductName, appName+strconv.FormatUint(o.ID, 10), "http://47.96.177.91:8888/api/order/verify", strconv.FormatFloat(o.Amount, 'f', 2, 64))
	if err != nil {
		t.Error("生成支付链接失败\t" + err.Error())
		models.TransactionRollback(trans)
		return
	}
	models.TransactionCommit(trans)

	t.Log(uri)
}
