package controllers

import (
	"MoShow/models"
	"MoShow/utils"

	"github.com/astaxie/beego"
)

//ConfigController .
type ConfigController struct {
	beego.Controller
}

//GetGiftList .
// @Title 获取礼物列表
// @Description 获取礼物列表
// @Success 200 {object} utils.ResultDTO
// @router /gifts [get]
func (c *ConfigController) GetGiftList() {
	dto := utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	conf := &models.Config{}
	val, err := conf.GetCommonGiftInfo()
	if err != nil {
		dto.Message = "获取礼物列表失败\t" + err.Error()
		beego.Error(err)
		return
	}

	dto.Data = val
	dto.Sucess = true
}
