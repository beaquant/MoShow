package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"time"

	"github.com/astaxie/beego"
)

//TimelineController .
type TimelineController struct {
	beego.Controller
}

//TimelineUserInfo 首页专区用户信息
type TimelineUserInfo struct {
	models.TimelineUser
	CoverInfo *models.UserCoverInfo `json:"cover_info"`
}

//NewCommers .
// @Title 新用户专区
// @Description 新用户专区,注册时间最近15天内的用户在此参与排序。(性别为1时，查询所有男用户,为0时，只查女主播，不查普通女用户)
// @Param   gender		    query    int  	false       "性别,0女,1男"
// @Param   skip		    query    int  	false       "偏移量"
// @Param   limit		    query    int  	false       "返回结果数限制"
// @Success 200 {object} models.UserProfile
// @router /newcomer [get]
func (c *TimelineController) NewCommers() {
	dto := utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	skip, err := c.GetInt("skip")
	if err != nil || skip < 0 {
		skip = 0
	}

	limit, err := c.GetInt("limit")
	if err != nil || limit < 1 {
		limit = 10
	}

	gender, err := c.GetInt("gender")
	if err != nil || gender < 0 || gender > 1 {
		gender = 0
	}

	ul, err := (&models.TimelineUser{}).QueryRecent(time.Now().AddDate(0, 0, -15).Unix(), gender, skip, limit)
	if err != nil {
		dto.Message = "获取新用户列表失败\t" + err.Error()
		beego.Error(err)
		return
	}

	var ti []TimelineUserInfo
	for index := range ul {
		ti = append(ti, TimelineUserInfo{TimelineUser: ul[index], CoverInfo: ul[index].GetCover()})
	}

	dto.Data = ti
	dto.Sucess = true
}

//Active .
// @Title 活跃专区
// @Description 活跃专区,所有用户参与查询。(性别为1时，查询所有男用户,为0时，只查女主播，不查普通女用户)
// @Param   gender		    query    int  	false       "性别,0女,1男"
// @Param   skip		    query    int  	false       "偏移量"
// @Param   limit		    query    int  	false       "返回结果数限制"
// @Success 200 {object} models.UserProfile
// @router /active [get]
func (c *TimelineController) Active() {
	dto := utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	skip, err := c.GetInt("skip")
	if err != nil || skip < 0 {
		skip = 0
	}

	limit, err := c.GetInt("limit")
	if err != nil || limit < 1 {
		limit = 10
	}

	gender, err := c.GetInt("gender")
	if err != nil || gender < 0 || gender > 1 {
		gender = 0
	}

	ul, err := (&models.TimelineUser{}).QueryAll(gender, skip, limit)
	if err != nil {
		dto.Message = "获取活跃用户列表失败\t" + err.Error()
		beego.Error(err)
		return
	}

	var ti []TimelineUserInfo
	for index := range ul {
		ti = append(ti, TimelineUserInfo{TimelineUser: ul[index], CoverInfo: ul[index].GetCover()})
	}

	dto.Data = ti
	dto.Sucess = true
}

//Suggestion .
// @Title 推荐专区
// @Description 推荐专区,用户由运营人员推荐
// @Param   skip		    query    int  	false       "偏移量"
// @Param   limit		    query    int  	false       "返回结果数限制"
// @Success 200 {object} models.UserProfile
// @router /suggestion [get]
func (c *TimelineController) Suggestion() {
	dto := utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	skip, err := c.GetInt("skip")
	if err != nil || skip < 0 {
		skip = 0
	}

	limit, err := c.GetInt("limit")
	if err != nil || limit < 1 {
		limit = 10
	}

	ul, err := (&models.TimelineUser{}).QueryHot(skip, limit)
	if err != nil {
		dto.Message = "获取推荐用户列表失败\t" + err.Error()
		beego.Error(err)
		return
	}

	var ti []TimelineUserInfo
	for index := range ul {
		ti = append(ti, TimelineUserInfo{TimelineUser: ul[index], CoverInfo: ul[index].GetCover()})
	}

	dto.Data = ti
	dto.Sucess = true
}
