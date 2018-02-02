package controllers

import (
	"MoShow/utils"
	"strconv"

	netease "github.com/MrSong0607/netease-im"
	"github.com/astaxie/beego"
)

//ImController 网易云信相关接口
type ImController struct {
	beego.Controller
}

//CreateImUser .
// @Title 创建网易云信用户
// @Description 创建网易云信用户
// @Param   name     		formData    string  	false       "网易云通信ID昵称，最大长度64字符，用来PUSH推送时显示的昵称"
// @Param   icon		    formData    string  	false       "头像"
// @Param   gender		    formData    int	  		false       "用户性别，0表示未知，1表示男，2女表示女"
// @Success 200 {object} utils.ResultDTO
// @router /create [put]
func (c *ImController) CreateImUser() {
	dto := utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	tk := GetToken(c.Ctx)
	user := &netease.ImUser{ID: strconv.FormatUint(tk.ID, 10)}

	if gender := c.GetString("gender"); len(gender) > 0 {
		if gd, err := strconv.Atoi(gender); err == nil && gd == 0 || gd == 1 || gd == 2 {
			user.Gender = gd
		}
	}

	if name := c.GetString("name"); len(name) > 0 {
		user.Name = name
	}

	if icon := c.GetString("icon"); len(icon) > 0 {
		user.IconURL = icon
	}

	imtk, err := utils.ImCreateUser(user)
	if err != nil {
		dto.Message = err.Error()
		return
	}

	dto.Data = imtk
	dto.Sucess = true
}

//RefreshToken .
// @Title 获取新token
// @Description 获取新token
// @Success 200 {object} utils.ResultDTO
// @router /refreshtoken [get]
func (c *ImController) RefreshToken() {
	dto := utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	tk := GetToken(c.Ctx)

	imtk, err := utils.ImRefreshToken(strconv.FormatUint(tk.ID, 10))
	if err != nil {
		dto.Message = err.Error()
		return
	}

	dto.Data = imtk
	dto.Sucess = true
}
