package controllers

import (
	"MoShow/utils"

	"github.com/astaxie/beego"
)

//UserController 用户节点
type UserController struct {
	beego.Controller
}

//Read .
// @Title 读取用户
// @Description 读取用户
// @Param   userid     path    string  true        "用户id,填me表示获取当前账号的用户信息"
// @Success 200 {object} utils.ResultDTO
// @router /:userid [get]
func (c *UserController) Read() {
	dto := utils.ResultDTO{}

	c.Data["json"] = dto
	c.ServeJSON()
}

//Update .
// @Title 更新用户
// @Description 更新用户
// @Param   userid     path    string  true        "The email for login"
// @Success 200 {object} utils.ResultDTO
// @router /:userid [post]
func (c *UserController) Update() {

}

//Del .
// @Title 删除用户
// @Description 删除用户
// @Param   userid     path    string  true        "The email for login"
// @Success 200 {object} utils.ResultDTO
// @router /:userid [delete]
func (c *UserController) Del() {

}

//SendGift .
// @Title 赠送礼物
// @Description 赠送礼物
// @Param   userid     path    		string  true        "赠送礼物的目标"
// @Param   giftid	   formData     uint     true		 "礼物id"
// @Param   count	   formData     uint     true		 "数量"
// @Success 200 {object} utils.ResultDTO
// @router /:userid [post]
func (c *UserController) SendGift() {

}
