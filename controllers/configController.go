package controllers

import (
	"MoShow/models"
	"MoShow/utils"

	"github.com/astaxie/beego"
)

//ConfigController 获取礼物列表，系统设置等
type ConfigController struct {
	beego.Controller
}

//GetCommonConfig .
// @Title 获取通用配置
// @Description 获取通用配置
// @Success 200 {object} utils.ResultDTO
// @router /common [get]
func (c *ConfigController) GetCommonConfig() {
	dto := utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	conf := &models.Config{}
	val, err := conf.GetCommonConfig()
	if err != nil {
		dto.Message = "获取通用配置失败\t" + err.Error()
		beego.Error(err, c.Ctx.Request.UserAgent())
		return
	}

	dto.Data = val
	dto.Sucess = true
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

//GetProductList .
// @Title 获取商品列表
// @Description 获取商品列表
// @Success 200 {object} utils.ResultDTO
// @router /products [get]
func (c *ConfigController) GetProductList() {
	dto := utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	conf := &models.Config{}
	val, err := conf.GetProductInfo()
	if err != nil {
		dto.Message = "获取商品列表失败\t" + err.Error()
		beego.Error(err, c.Ctx.Request.UserAgent())
		return
	}

	dto.Data = val
	dto.Sucess = true
}

//GetCosSign .
// @Title 获取对象存储签名
// @Description 获取对象存储签名
// @Param   dir			query    string  	true       "文件目录"
// @Success 200 {object} utils.ResultDTO
// @router /cossign [get]
func (c *ConfigController) GetCosSign() {
	dto := utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	dir := c.GetString("dir")

	sign, err := utils.GetTecentImgSignV5(dir)
	if err != nil {
		beego.Error("获取腾讯COS签名失败\t"+err.Error(), c.Ctx.Request.UserAgent())
		dto.Message = "获取腾讯COS签名失败\t" + err.Error()
		return
	}

	dto.Data = sign
	dto.Sucess = true
}
