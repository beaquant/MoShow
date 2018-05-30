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
	UserProfileInfo
	CreatedAt int64 `json:"create_at" gorm:"column:create_at"`
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

	cate := c.GetString("cate")
	if len(cate) == 0 {
		cate = "suggestion"
	}

	ul := []models.TimelineUser{}
	qs := models.GetContext().Table((models.UserProfile{}).TableName()).Select("user_profile.*,create_at").
		Joins("left join users on users.id = user_profile.id").Where("users.id <> 1 and user_profile.user_status <> ? and users.acct_status <> ?", models.UserStatusBlock, models.AcctStatusShield)
	if up.UserType == models.UserTypeFaker {
		qs = qs.Where("user_type = ?", models.UserTypeFaker).Order("dial_duration desc")
	} else {
		if up.UserType != models.UserTypeAnchor {
			qs = qs.Where("user_type = ?", models.UserTypeAnchor).Where("gender = ?", models.GenderWoman)
		} else {
			qs = qs.Where("user_type <> ?", models.UserTypeFaker).Where("gender = ?", models.GenderMan)
		}
	}

	switch cate {
	case "newcomer":
		if up.UserType != models.UserTypeFaker {
			qs = qs.Where("create_at > ?", time.Now().AddDate(0, 0, -15).Unix())
		}
		err = qs.Order("online_status = 1 or online_status = 2 desc, id desc").Offset(skip).Limit(limit).Scan(&ul).Error
	case "active":
		err = qs.Order("online_status = 1 or online_status = 2 desc, recent_duration desc").Offset(skip).Limit(limit).Scan(&ul).Error
	case "suggestion":
		if up.UserType != models.UserTypeFaker {
			qs = qs.Where("user_status = ?", models.UserStatusHot)
		}
		err = qs.Order("online_status = 1 or online_status = 2 desc, recent_duration desc").Offset(skip).Limit(limit).Scan(&ul).Error
	}

	if err != nil {
		dto.Message = "获取用户列表失败\t" + err.Error()
		beego.Error(err)
		return
	}

	var ti []TimelineUserInfo
	for index := range ul {
		upi := &UserProfileInfo{UserProfile: ul[index].UserProfile}
		genUserPorfileInfoCommon(upi, upi.GetCover())
		ti = append(ti, TimelineUserInfo{UserProfileInfo: *upi, CreatedAt: ul[index].CreatedAt})
	}

	tli := TimelineInfo{Users: ti}
	if config, _ := (&models.Config{}).GetCommonConfig(); config != nil {
		if up.UserType == models.UserTypeFaker {
			tli.Banners = config.CheckModeBanners
		} else {
			tli.Banners = config.Banners
		}
	}
	dto.Data = tli
	dto.Sucess = true
}
