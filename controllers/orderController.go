package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"strconv"

	"github.com/astaxie/beego"
)

//OrderController 充值，支付，提现等接口
type OrderController struct {
	beego.Controller
}

//Detail  订单详情
// @Title 订单详情
// @Description 订单详情
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
