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

//Create .
// @Title 创建用户信息
// @Description 创建用户信息
// @Param   alias     		formData    string  	true        "昵称"
// @Param   gender    		formData    int     	true        "性别，男(1),女(0)"
// @Param   cover_pic	    formData    string  	true        "照片"
// @Param   description     formData    string  	true        "签名"
// @Param   birthday     	formData    time.Time  	true        "生日"
// @Param   location     	formData    string  	true        "地区"
// @Param   price	     	formData    uint	  	false       "价格"
// @Success 200 {object} 	utils.ResultDTO
// @router /create [put]
func (c *UserController) Create() {
	tk := GetToken(c.Ctx)
	dto := &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	up := &models.UserProfile{ID: tk.ID}
	if alias := c.GetString("alias"); len(alias) > 0 {
		up.Alias = alias
	}

	err := up.Add()
	if err != nil {
		beego.Error(err)
		dto.Message = err.Error()
	}
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
