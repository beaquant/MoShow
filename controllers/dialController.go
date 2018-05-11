package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"net/http"
	"strconv"

	netease "github.com/MrSong0607/netease-im"
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
	Parter UserProfileInfo
}

func genDialInfo(self uint64, dials []models.Dial) []DialInfo {
	var dis []DialInfo
	upkv := make(map[uint64]*UserProfileInfo)

	for index := range dials {
		di := DialInfo{Dial: dials[index]}
		u := &UserProfileInfo{}
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

	if id == 0 {
		dto.Message = "通话记录ID不能为0"
		dto.Code = utils.DtoStatusParamError
		return
	}

	dial := &models.Dial{ID: id}
	if err := dial.Read(); err != nil {
		beego.Error("获取通话记录失败", err, c.Ctx.Request.UserAgent(), "通话记录ID:", id)
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

//NmCallback .
// @router /nmcallback [post]
func (c *DialController) NmCallback() {
	bd, err := utils.ImClient.GetEventNotification(c.Ctx.Request)
	if err != nil {
		beego.Error("云信抄送异常", err)
		c.Abort(strconv.Itoa(http.StatusBadRequest))
		return
	}

	kv := make(map[string]interface{})
	if err := utils.JSONUnMarshalFromByte(bd, &kv); err != nil {
		beego.Error("云信回执解析异常", err, "body", string(bd))
		c.Abort(strconv.Itoa(http.StatusBadRequest))
		return
	}

	val, ok := kv["eventType"]
	if !ok {
		beego.Error("云信回执内容错误", string(bd))
		c.Ctx.Output.Body([]byte("success"))
		return
	}

	v, ok := val.(string)
	if !ok {
		c.Abort(strconv.Itoa(http.StatusBadRequest))
		return
	}

	switch v {
	case netease.EventTypeMediaDuration:
		ci := &netease.AudioCopyInfo{}
		if err := utils.JSONUnMarshalFromByte(bd, ci); err != nil {
			beego.Error("云信回执解析错误", err, "body", string(bd))
			c.Abort(strconv.Itoa(http.StatusBadRequest))
			return
		}

		dl := &models.Dial{}
		if err := dl.ReadFromNimID(ci.ChannelID); err != nil { //找不到通话记录，丢弃该回执
			beego.Error("云信回执找不到指定的聊天通道相关的通话记录", err, "body", string(bd))
			break
		}

		if err := dl.UpdateNmAudioCopy(ci); err != nil {
			beego.Error("更新云信通话时长信息回执失败", err)
			c.Abort(strconv.Itoa(http.StatusBadRequest))
			return
		}
	case netease.EventTypeMediaInfo:
		dci := &netease.AudioDownloadCopyInfo{}
		if err := utils.JSONUnMarshalFromByte(bd, dci); err != nil {
			beego.Error("云信回执解析错误", err, "body", string(bd))
			c.Abort(strconv.Itoa(http.StatusBadRequest))
			return
		}

		var fi []netease.FileDownloadInfo
		if err := utils.JSONUnMarshal(dci.FileInfo, &fi); err != nil {
			beego.Error("云信回执解析错误", err, "body", string(bd))
			c.Abort(strconv.Itoa(http.StatusBadRequest))
			return
		}

		if fi == nil || len(fi) == 0 { //文件下载信息为空，丢弃该回执
			break
		}

		dl := &models.Dial{}
		if err := dl.ReadFromNimID(fi[0].ChannelID); err != nil { //找不到通话记录，丢弃该回执
			beego.Error("云信回执找不到指定的聊天通道相关的通话记录", err, "body", string(bd))
			break
		}

		if err := dl.UpdateNmAudioDlCopy(dci, fi); err != nil {
			beego.Error("更新云信视频录制地址信息失败", err)
			c.Abort(strconv.Itoa(http.StatusBadRequest))
			return
		}
	}

	c.Ctx.Output.Body([]byte("success"))
}
