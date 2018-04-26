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

//TimelineInfo 首页专区用户信息
type TimelineInfo struct {
	Users   []TimelineUserInfo `json:"users"`
	Banners []models.Banner    `json:"banners"`
}

//TimelineUserInfo .
type TimelineUserInfo struct {
	UserPorfileInfo
	CreatedAt int64  `json:"create_at" gorm:"column:create_at"`
	Duration  uint64 `json:"recent_duration" gorm:"column:recent_duration"`
}

//Users .
// @Title 首页专区列表
// @Description newcomer:新用户专区,注册时间最近15天内的用户在此参与排序。active:活跃专区,所有用户参与查询。 suggestion:推荐专区,用户由运营人员推荐
// @Param   cate		    query    string  	true       	"专区分类,newcomer,active,suggestion"
// @Param   skip		    query    int  		false       "偏移量"
// @Param   limit		    query    int  		false       "返回结果数限制"
// @Success 200 {object} models.UserProfile
// @router /users [get]
func (c *TimelineController) Users() {
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
	var faker bool
	if up.Gender == models.GenderWoman && up.UserType == models.UserTypeAnchor {
		gender = models.GenderMan
	} else {
		gender = models.GenderWoman
	}

	if up.UserType == models.UserTypeFaker {
		faker = true
	}

	cate := c.GetString("cate")
	if len(cate) == 0 {
		cate = "suggestion"
	}

	var ul []models.TimelineUser
	switch cate {
	case "newcomer":
		ul, err = (&models.TimelineUser{}).QueryRecent(faker, time.Now().AddDate(0, 0, -15).Unix(), gender, skip, limit)
	case "active":
		ul, err = (&models.TimelineUser{}).QueryAll(faker, gender, skip, limit)
	case "suggestion":
		ul, err = (&models.TimelineUser{}).QuerySuggestion(faker, gender, skip, limit)
	}

	if err != nil {
		dto.Message = "获取用户列表失败\t" + err.Error()
		beego.Error(err)
		return
	}

	var ti []TimelineUserInfo
	for index := range ul {
		upi := &UserPorfileInfo{UserProfile: ul[index].UserProfile}
		genUserPorfileInfoCommon(upi, upi.GetCover())
		ti = append(ti, TimelineUserInfo{UserPorfileInfo: *upi, CreatedAt: ul[index].CreatedAt, Duration: ul[index].Duration})
	}

	tli := TimelineInfo{Users: ti}
	if config, _ := (&models.Config{}).GetCommonConfig(); config != nil {
		tli.Banners = config.Banners
	}
	dto.Data = tli
	dto.Sucess = true
}
