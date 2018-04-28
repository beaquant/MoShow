package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"strconv"

	"github.com/astaxie/beego"
)

//BalanceChgController 账户余额相关接口
type BalanceChgController struct {
	beego.Controller
}

//BalanceChgInfo .
type BalanceChgInfo struct {
	models.BalanceChg
	FromUserInfo *UserProfileInfo `json:"from_user_info"`
}

//GetIncomeList .
// @Title 获取收益列表
// @Description 获取收益列表
// @Param   length     	query    int  	true       "长度"
// @Param   skip		query    int  	true       "偏移量"
// @Success 200 {object} utils.ResultDTO
// @router /incomes [get]
func (c *BalanceChgController) GetIncomeList() {
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

	chg := &models.BalanceChg{UserID: tk.ID}
	lst, err := chg.GetIncomeChgs(len, skip)
	if err != nil {
		beego.Error("查询变动失败:\t"+err.Error(), c.Ctx.Request.UserAgent())
		dto.Message = "查询变动失败:\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	dto.Data = lst
	dto.Sucess = true
}

//GetPaymentList .
// @Title 获取支出列表
// @Description 获取支出列表
// @Param   length     	query    int  	true       "长度"
// @Param   skip		query    int  	true       "偏移量"
// @Success 200 {object} utils.ResultDTO
// @router /payments [get]
func (c *BalanceChgController) GetPaymentList() {
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

	chg := &models.BalanceChg{UserID: tk.ID}
	lst, err := chg.GetPaymentChgs(len, skip)
	if err != nil {
		beego.Error("查询变动失败:\t"+err.Error(), c.Ctx.Request.UserAgent())
		dto.Message = "查询变动失败:\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	dto.Data = lst
	dto.Sucess = true
}

//GetInviteIncomList .
// @Title 获取邀请收入列表
// @Description 获取邀请收入列表
// @Param   length     	query    int  	true       "长度"
// @Param   skip		query    int  	true       "偏移量"
// @Success 200 {object} utils.ResultDTO
// @router /ivtincome [get]
func (c *BalanceChgController) GetInviteIncomList() {
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

	chg := &models.BalanceChg{UserID: tk.ID}
	lst, err := chg.GetInviteIncomeChgs(len, skip)
	if err != nil {
		beego.Error("查询变动失败:\t"+err.Error(), c.Ctx.Request.UserAgent())
		dto.Message = "查询变动失败:\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	var bcis []BalanceChgInfo
	for index := range lst {
		bci := BalanceChgInfo{BalanceChg: lst[index], FromUserInfo: &UserProfileInfo{UserProfile: models.UserProfile{ID: lst[index].FromUserID}}}
		if bci.FromUserInfo.ID != 0 {
			if err := bci.FromUserInfo.Read(); err != nil {
				beego.Error("获取用户信息出错", err)
			}
		}

		genUserPorfileInfoCommon(bci.FromUserInfo, bci.FromUserInfo.GetCover())
		bcis = append(bcis, bci)
	}

	dto.Data = bcis
	dto.Sucess = true
}

//GetChgDetail .
// @Title 获取单条变动详情
// @Description 获取单条变动详情
// @Param   chgid     		path    	int  	true        "变动id"
// @Success 200 {object} utils.ResultDTO
// @router /:chgid [get]
func (c *BalanceChgController) GetChgDetail() {
	dto, chgidStr := utils.ResultDTO{}, c.Ctx.Input.Param(":chgid")
	defer dto.JSONResult(&c.Controller)

	chgid, err := strconv.ParseUint(chgidStr, 10, 64)
	if err != nil {
		beego.Error("参数解析错误:chgid\t"+err.Error(), c.Ctx.Request.UserAgent(), chgidStr)
		dto.Message = "参数解析错误:chgid\t" + err.Error()
		dto.Code = utils.DtoStatusParamError
		return
	}

	chg := &models.BalanceChg{ID: chgid}
	if err := chg.Read(); err != nil {
		beego.Error("查询变动失败:\t"+err.Error(), c.Ctx.Request.UserAgent(), chgid)
		dto.Message = "查询变动失败:\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	dto.Data = chg
	dto.Sucess = true
}
