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

//DialDetail 通话详情
type DialDetail struct {
	models.Dial
	models.ClearingInfo
}

//DialInfo .
type DialInfo struct {
	Dial   models.Dial
	Parter UserPorfileInfo
}

func genDialInfo(self uint64, dials []models.Dial) []DialInfo {
	var dis []DialInfo
	upkv := make(map[uint64]*UserPorfileInfo)

	for index := range dials {
		di := DialInfo{Dial: dials[index]}
		u := &UserPorfileInfo{}
		if dials[index].FromUserID == self {
			u.ID = dials[index].ToUserID
		} else {
			u.ID = dials[index].FromUserID
		}

		if upi, ok := upkv[u.ID]; ok {
			di.Parter = *upi
		} else {
			u.Read()
			genUserPorfileInfoCommon(u, u.GetCover())
			upkv[u.ID] = u
			di.Parter = *u
		}

		dis = append(dis, di)
	}

	return dis
}

func genDialDetail(d *models.Dial) (*DialDetail, error) {
	ci := &models.ClearingInfo{}
	var err error

	if len(d.Clearing) > 0 {
		err = utils.JSONUnMarshal(d.Clearing, ci)

	}
	return &DialDetail{Dial: *d, ClearingInfo: *ci}, err
}

//DialList .
// @Title 获取通话记录列表
// @Description 获取通话记录列表
// @Param   length     	query    int  	true       "长度,最大20"
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

	if len > 20 {
		len = 20
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

	dto.Data = genDialInfo(tk.ID, lst)
	dto.Sucess = true
}

//GetDialDetail 通话记录详情
// @Title 通话记录详情
// @Description 通话记录详情
// @Param   dialid     	path    int  	true       "通话记录ID"
// @Success 200 {object} utils.ResultDTO
// @router /:dialid [get]
func (c *DialController) GetDialDetail() {
	dto, tk := utils.ResultDTO{}, GetToken(c.Ctx)
	defer dto.JSONResult(&c.Controller)

	id, err := strconv.ParseUint(c.Ctx.Input.Param(":dialid"), 10, 64)
	if err != nil {
		beego.Error("参数解析错误:dialid\t", err, c.Ctx.Request.UserAgent(), id)
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

	if dial.FromUserID != tk.ID && dial.ToUserID != tk.ID {
		beego.Error("当前用户没有此通话记录")
		dto.Message = "当前用户没有此通话记录"
		dto.Code = utils.DtoStatusParamError
		return
	}

	dd, err := genDialDetail(dial)
	if err != nil {
		beego.Error("获取通话记录", err)
		dto.Message = "获取通话记录失败" + err.Error()
		return
	}

	dto.Data = dd
	dto.Sucess = true
	dto.Message = "获取成功"
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
