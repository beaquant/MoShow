package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

//DialController .
type DialController struct {
	beego.Controller
}

//DialList .
// @Title 获取通话记录列表
// @Description 获取通话记录列表
// @Param   length     	query    int  	true       "长度"
// @Param   skip		query    int  	true       "偏移量"
// @Success 200 {object} utils.ResultDTO
// @router /list [get]
func (c *DialController) DialList() {
	dto, tk := utils.ResultDTO{}, GetToken(c.Ctx)
	defer dto.JSONResult(&c.Controller)

	len, err := c.GetInt("length")
	if err != nil {
		beego.Error("参数解析错误:length\t"+err.Error(), c.Ctx.Request.UserAgent(), c.GetString("length"))
		dto.Message = "参数解析错误:length\t" + err.Error()
		dto.Code = utils.DtoStatusParamError
		return
	}

	skip, err := c.GetInt("skip")
	if err != nil {
		beego.Error("参数解析错误:skip\t"+err.Error(), c.Ctx.Request.UserAgent(), c.GetString("skip"))
		dto.Message = "参数解析错误:skip\t" + err.Error()
		dto.Code = utils.DtoStatusParamError
		return
	}

	lst, err := (&models.Dial{}).GetDialList(tk.ID, len, skip)
	if err != nil {
		beego.Error("查询通话记录列表失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "查询通话记录列表失败\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	dto.Data = lst
	dto.Sucess = true
}

//Del .
// @Title 删除通话记录
// @Description 删除通话记录
// @Param   dialid     	path    int  	true       "通话记录id"
// @Success 200 {object} utils.ResultDTO
// @router /:dialid [delete]
func (c *DialController) Del() {
	dto, tk := utils.ResultDTO{}, GetToken(c.Ctx)
	defer dto.JSONResult(&c.Controller)

	ids := strings.TrimSpace(c.Ctx.Input.Param(":dialid"))
	id, err := strconv.ParseUint(ids, 10, 64)
	if err != nil {
		beego.Error("参数解析错误:dialid\t", err, c.Ctx.Request.UserAgent(), ids)
		dto.Message = "参数解析错误:dialid\t" + err.Error()
		dto.Code = utils.DtoStatusParamError
		return
	}

	dial := &models.Dial{ID: id}
	if err := dial.Read(); err != nil {
		beego.Error("获取通话记录失败", err, c.Ctx.Request.UserAgent(), id)
		dto.Message = "获取通话记录失败" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	if dial.FromUserID != tk.ID {
		beego.Error("不能删除他人的通话记录")
		dto.Message = "该通话记录不属于当前用户"
		dto.Code = utils.DtoStatusParamError
		return
	}

	if err := dial.Del(); err != nil {
		beego.Error("删除通话记录失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "删除通话记录失败" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	dto.Message = "删除成功"
	dto.Sucess = true
}
