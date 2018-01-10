package controllers

import (
	"MoShow/utils"

	"github.com/astaxie/beego"
)

//UserController 用户节点
type UserController struct {
	beego.Controller
}

//Create .
// @Title 创建用户
// @Description 创建用户
// @Success 200 {object} utils.ResultDTO
// @router /create [put]
func (c *UserController) Create() {
	dto := utils.ResultDTO{}

	c.Data["json"] = dto
	c.ServeJSON()
}

//Read .
// @Title 读取用户
// @Description 读取用户
// @Param   userid     path    string  true        "用户id,填me表示获取当前账号的用户信息"
// @Success 200 {object} utils.ResultDTO
// @router /:userid [get]
func (c *UserController) Read() {

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
