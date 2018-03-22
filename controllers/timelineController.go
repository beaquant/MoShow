package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"time"

	"github.com/astaxie/beego"
)

//TimelineController 首页专区推荐相关接口
type TimelineController struct {
	beego.Controller
}

//TimelineUserInfo 首页专区用户信息
type TimelineUserInfo struct {
	models.TimelineUser
	CoverInfo *models.UserCoverInfo `json:"cover_info"`
}

//NewCommers .
// @Title 首页专区列表
// @Description newcomer:新用户专区,注册时间最近15天内的用户在此参与排序。active:活跃专区,所有用户参与查询。 suggestion:推荐专区,用户由运营人员推荐
// @Param   cate		    query    string  	true       	"专区分类,newcomer,active,suggestion"
// @Param   skip		    query    int  		false       "偏移量"
// @Param   limit		    query    int  		false       "返回结果数限制"
// @Success 200 {object} models.UserProfile
// @router /users [get]
func (c *TimelineController) NewCommers() {
	dto, tk := utils.ResultDTO{}, GetToken(c.Ctx)
	defer dto.JSONResult(&c.Controller)

	skip, err := c.GetInt("skip")
	if err != nil || skip < 0 {
		skip = 0
	}

	limit, err := c.GetInt("limit")
	if err != nil || limit < 1 {
		limit = 10
	}

	up := &models.UserProfile{ID: tk.ID}
	if err := up.Read(); err != nil {
		beego.Error("获取用户信息失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "获取用户信息失败\t" + err.Error()
		return
	}

	var gender int
	if up.Gender == models.GenderMan {
		gender = models.GenderWoman
	} else if up.Gender == models.GenderWoman {
		gender = models.GenderMan
	} else {
		gender = models.GenderDefault
	}

	cate := c.GetString("cate")
	if len(cate) == 0 {
		cate = "suggestion"
	}

	var ul []models.TimelineUser
	switch cate {
	case "newcomer":
		ul, err = (&models.TimelineUser{}).QueryRecent(time.Now().AddDate(0, 0, -15).Unix(), gender, skip, limit)
	case "active":
		ul, err = (&models.TimelineUser{}).QueryAll(gender, skip, limit)
	case "suggestion":
		ul, err = (&models.TimelineUser{}).QueryHot(skip, limit)
	}

	if err != nil {
		dto.Message = "获取新用户列表失败\t" + err.Error()
		beego.Error(err)
		return
	}

	var ti []*UserPorfileInfo
	for index := range ul {
		upi := &UserPorfileInfo{UserProfile: ul[index].UserProfile}
		genUserPorfileInfoCommon(upi, upi.GetCover())
		ti = append(ti, upi)
	}

	dto.Data = ti
	dto.Sucess = true
}
