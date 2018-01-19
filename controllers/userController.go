package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"strconv"
	"strings"

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
	dto := utils.ResultDTO{Sucess: false}
	defer dto.JSONResult(&c.Controller)

	var uid uint64
	var err error
	tk := GetToken(c.Ctx)
	uidStr := strings.TrimSpace(c.Ctx.Input.Param(":userid"))
	if uidStr == "me" {
		uid = tk.ID
	} else {
		uid, err = strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			beego.Error(err)
			dto.Message = err.Error()
			return
		}
	}

	up := &models.UserProfile{ID: uid}
	err = up.Read()
	if err != nil {
		beego.Error(err)
		dto.Message = err.Error()
		return
	}

	dto.Data = up
	dto.Sucess = true
}

//Update .
// @Title 更新用户
// @Description 更新用户
// @Success 200 {object} utils.ResultDTO
// @router /update [post]
func (c *UserController) Update() {
	dto := utils.ResultDTO{Sucess: false}
	defer dto.JSONResult(&c.Controller)
	dto.Sucess = true
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
// @router /:userid/sendgift [post]
func (c *UserController) SendGift() {

}
