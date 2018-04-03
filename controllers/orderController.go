package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

const appName = "MoShow"

var re = regexp.MustCompile("\\d+")

//OrderController 充值，支付，提现等接口
type OrderController struct {
	beego.Controller
}

//Detail  订单详情
// @Title 订单详情
// @Description 订单详情
// @Param   orderid     		path    	string  	true        "订单id"
// @Success 200 {object} utils.ResultDTO
// @router /:orderid/detail [get]
func (c *OrderController) Detail() {
	dto, oidStr := &utils.ResultDTO{}, c.Ctx.Input.Param(":userid")
	defer dto.JSONResult(&c.Controller)

	orderID, err := strconv.ParseUint(oidStr, 10, 64)
	if err != nil {
		beego.Error(err)
		dto.Message = "订单ID格式错误\t" + err.Error()
		return
	}

	order := &models.Order{ID: orderID}
	if order.Read() != nil {
		beego.Error(err, c.Ctx.Request.UserAgent())
		dto.Message = "订单查询失败\t" + err.Error()
		return
	}

	dto.Data = order
	dto.Sucess = true
}

//CreateWebPay .
// @Title 创建网页订单
// @Description 创建网页订单
// @Param   prodid		    formData    int  	true       "套餐ID"
// @Success 200 {object} utils.ResultDTO
// @router /webpay [post]
func (c *OrderController) CreateWebPay() {
	dto, tk := &utils.ResultDTO{}, GetToken(c.Ctx)
	defer dto.JSONResult(&c.Controller)

	pid, err := c.GetUint64("prodid")
	if err != nil {
		beego.Error("产品ID格式错误", err, c.Ctx.Request.UserAgent())
		dto.Message = "产品ID格式错误\t" + err.Error()
		return
	}

	prod, err := getProdFromID(pid)
	if err != nil {
		beego.Error("产品ID错误", pid, err, c.Ctx.Request.UserAgent())
		dto.Message = "产品ID错误\t" + err.Error()
		return
	}

	prodInfo, err := utils.JSONMarshalToString(prod)
	if err != nil {
		beego.Error("解析产品信息出错")
		dto.Message = "解析产品信息出错"
		return
	}

	trans := models.TransactionGen()
	o := &models.Order{Amount: prod.Price, CoinCount: prod.CoinCount + prod.Extra, UserID: tk.ID, PayType: models.PayTypeAlipay, CreateAt: time.Now().Unix(), ProductInfo: prodInfo, PayInfo: "{}"}
	if err := o.Add(trans); err != nil {
		beego.Error("添加订单失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "添加订单失败\t" + err.Error()
		models.TransactionRollback(trans)
		return
	}

	//用APPName加上订单ID拼接成唯一交易号(支付宝规定每个收款账号下面的交易号必须唯一)
	uri, err := utils.CreatePayment(prod.ProductName, appName+strconv.FormatUint(o.ID, 10), "http://47.96.177.91:8888/api/order/verify", strconv.FormatFloat(o.Amount, 'f', 2, 64))
	if err != nil {
		beego.Error("生成支付链接失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "生成支付链接失败\t" + err.Error()
		models.TransactionRollback(trans)
		return
	}
	models.TransactionCommit(trans)

	beego.Info(uri)
	dto.Data = uri
	dto.Sucess = true
}

//AlipayConfirm .
//@Title 支付宝回调
// @Description 支付宝回调
// @Param	userinfo		body 	common.ResultDTO	true		"测试用"
//@router /verify [post]
func (c *OrderController) AlipayConfirm() {
	noti, err := utils.ConfirmPayment(c.Ctx.Request)
	if beego.BConfig.RunMode == "dev" {
		b, _ := json.Marshal(noti)
		beego.Debug(string(b))
	}

	if err != nil {
		beego.Error("支付宝订单回调验证失败", err, noti)
		return
	}

	//正则表达式解析订单ID
	orderID, err := strconv.ParseUint(re.FindString(noti.OutTradeNo), 10, 64)
	if err != nil {
		beego.Error("解析支付宝订单回调参数失败", err)
		return
	}

	amount, err := strconv.ParseFloat(noti.TotalAmount, 64)
	if err != nil {
		beego.Error("解析支付宝订单回调参数失败", err)
		return
	}

	order := models.Order{ID: orderID}
	if err := order.Read(); err != nil {
		beego.Error("获取订单详情失败", err)
		return
	}

	if order.Success {
		beego.Error("该订单已确认", noti)
		c.Ctx.Output.Body([]byte("success"))
		return
	}

	if order.Amount != amount {
		beego.Error("支付宝订单金额异常", noti)
		return
	}

	up, param, chg, trans := &models.UserProfile{ID: order.UserID}, make(map[string]interface{}), &models.BalanceChg{UserID: order.UserID, Amount: int(order.CoinCount), Time: time.Now().Unix(), ChgType: models.BalanceChgTypeRecharge, ChgInfo: order.ProductInfo}, models.TransactionGen()
	param["pay_info"], _ = utils.JSONMarshalToString(noti)
	param["success"] = true
	param["pay_time"] = time.Now().Unix()

	if err := order.Update(param, trans); err != nil {
		beego.Error("更新订单状态失败", err)
		models.TransactionRollback(trans)
		return
	}

	if err := up.AddBalance(int(order.CoinCount), trans); err != nil {
		beego.Error("增加用户余额失败", err)
		models.TransactionRollback(trans)
		return
	}

	if err := chg.Add(trans); err != nil {
		beego.Error("用户充值，生成变动失败", err)
		models.TransactionRollback(trans)
		return
	}

	if err := (&models.UserExtra{ID: order.UserID}).AddBalanceHis(order.CoinCount, trans); err != nil {
		beego.Error("用户充值，增加历史余额失败", err)
		models.TransactionRollback(trans)
		return
	}

	models.TransactionCommit(trans)
	c.Ctx.Output.Body([]byte("success"))
}

func getProdFromID(pid uint64) (*models.Product, error) {
	plst, err := (&models.Config{}).GetProductInfo()
	if err != nil {
		return nil, err
	}

	for index := range plst {
		if plst[index].ID == pid {
			return &plst[index], nil
		}
	}

	return nil, errors.New("未能找到指定套餐")
}
