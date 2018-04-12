package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

//UserController 用户信息查询，更新等接口
type UserController struct {
	beego.Controller
}

//UserPorfileInfo 用户信息
type UserPorfileInfo struct {
	models.UserProfile
	ImTk        string                 `json:"im_token,omitempty"`
	Alipay      *models.AlipayAcctInfo `json:"alipay_acct,omitempty"`
	Followed    bool                   `json:"followed" description:"是否已关注"`
	IsFill      bool                   `json:"is_fill" description:"资料是否完善"`
	AnswerRate  float64                `json:"answer_rate" description:"接通率"`
	CheckStatus *models.ProfileChg     `json:"check_status" description:"审核状态"`
	Avatar      string                 `json:"avatar"`
	Gallery     []string               `json:"gallery"`
	Video       string                 `json:"video"`
	VideoPayed  bool                   `json:"video_payed"`
	GiftRecv    []models.GiftHisInfo   `json:"gift_recv"`
}

//UserOperateInfo .
type UserOperateInfo struct {
	User   *UserPorfileInfo `json:"user"`
	OpTime int64            `json:"time"`
}

//Read .
// @Title 读取用户
// @Description 读取用户
// @Param   userid     path    string  true        "用户id,填me表示获取当前账号的用户信息"
// @Success 200 {object} models.UserProfile
// @router /:userid [get]
func (c *UserController) Read() {
	dto := utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	var uid uint64
	var err error
	tk := GetToken(c.Ctx)
	uidStr := strings.TrimSpace(c.Ctx.Input.Param(":userid"))
	if uidStr == "me" {
		uid = tk.ID
	} else {
		uid, err = strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			beego.Error(err)
			dto.Message = err.Error()
			return
		}
	}

	up := &models.UserProfile{ID: uid}
	err = up.Read()
	if err != nil {
		beego.Error(err)
		dto.Message = err.Error()
		return
	}

	upi := &UserPorfileInfo{UserProfile: *up}
	if uid == tk.ID {
		if upi, err = genSelfUserPorfileInfo(up, nil); err != nil {
			beego.Error("获取用户资料失败", err, c.Ctx.Request.UserAgent())
			dto.Message = "获取用户资料失败" + err.Error()
			return
		}
	} else {
		//自己查看自己时不增加记录
		if err := (&models.Guest{}).AddView(uid, tk.ID); err != nil { //增加访客记录
			beego.Error(err)
			dto.Message = "增加访客记录失败\t" + err.Error()
		}
		genUserPorfileInfoCommon(upi, up.GetCover())
		if len(upi.Video) > 0 {
			upi.VideoPayed, err = (&models.BalanceChg{UserID: tk.ID}).IsVideoPayed(upi.Video, upi.ID)
			if err != nil {
				beego.Error("获取视频付费信息失败", err, c.Ctx.Request.UserAgent())
				dto.Message = "获取视频付费信息失败" + err.Error()
				dto.Code = utils.DtoStatusDatabaseError
				return
			}
		}
	}

	lst, err := (&models.UserExtra{ID: uid}).GetGiftHis()
	if err != nil {
		beego.Error("获取礼物历史记录失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "获取礼物历史记录失败\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}
	upi.GiftRecv = lst

	sb := &models.Subscribe{ID: uid}
	if err := sb.Read(nil); err != nil {
		beego.Error("获取订阅信息", err, c.Ctx.Request.UserAgent())
		dto.Message = "获取订阅信息" + err.Error()
		return
	}

	if fl := sb.GetFollowers(); fl != nil {
		if _, ok := fl[tk.ID]; ok {
			upi.Followed = true
		}
	}

	dto.Data = upi
	dto.Sucess = true
}

//ReadExtra .
// @Title 读取用户统计信息
// @Description 读取用户统计信息
// @Success 200 {object} models.UserExtra
// @router /extra [get]
func (c *UserController) ReadExtra() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	ue := &models.UserExtra{ID: tk.ID}
	if err := ue.Read(); err != nil {
		beego.Error("获取用户统计信息失败", err)
		dto.Message = "获取用户统计信息失败" + err.Error()
	}

	dto.Data = ue
	dto.Sucess = true
	dto.Message = "获取成功"
}

//Update .
// @Title 更新用户
// @Description 更新用户
// @Param   alias     		formData    string  	false        "昵称"
// @Param   cover_pic	    formData    string  	false       "头像"
// @Param   gender		    formData    int  		false       "性别,男:1,女:2(性别一经设置，不能再修改)"
// @Param   gallery		    formData    []string  	false       "相册"
// @Param   video		    formData    string  	false       "视频"
// @Param   description     formData    string  	false       "签名"
// @Param   birthday     	formData    time.Time  	false       "生日,格式:2006-01-02"
// @Param   location     	formData    string  	false       "地区"
// @Param   price	     	formData    uint	  	false       "价格"
// @Success 200 {object} models.UserProfile
// @router /update [post]
func (c *UserController) Update() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	up := &models.UserProfile{ID: tk.ID}
	if err := up.Read(); err != nil {
		beego.Error("获取用户资料失败", err, c.Ctx.Request.UserAgent())
		dto.Message = err.Error()
		return
	}

	pc := &models.ProfileChg{ID: tk.ID}
	if err := pc.Read(nil); err != nil {
		beego.Error("获取用户信息变动失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "获取用户信息变动失败" + err.Error()
		return
	}

	cv, imgChg, param, pcParam := up.GetCover(), false, make(map[string]interface{}), make(map[string]interface{})
	if cv == nil {
		cv = &models.UserCoverInfo{}
	}

	if alias := c.GetString("alias"); len(alias) > 0 {
		param["alias"], up.Alias = alias, alias
	}

	if up.Gender == models.GenderDefault {
		if genderStr := c.GetString("gender"); len(genderStr) > 0 {
			gender, err := strconv.Atoi(genderStr)
			if err != nil || (gender != 1 && gender != 2) {
				beego.Error("参数解析错误:"+genderStr, err, c.Ctx.Request.UserAgent())
				dto.Message = err.Error()
				dto.Code = utils.DtoStatusParamError
				return
			}

			if gender == 1 {
				param["gender"] = models.GenderMan
				up.Gender = models.GenderMan
			} else if gender == 2 {
				param["gender"] = models.GenderWoman
				up.Gender = models.GenderWoman
			}
		}
	}

	if coverPic := c.GetString("cover_pic"); len(coverPic) > 0 {
		cv.CoverPicture, imgChg = &models.Picture{ImageURL: coverPic}, false //头像先审核再更新
		pcParam["cover_pic"] = coverPic                                      //用户头像信息变动
		pcParam["cover_pic_check"] = models.CheckStatusUncheck               //用户头像信息变动状态
	}

	if gallery := c.GetStrings("gallery"); gallery != nil && len(gallery) > 0 {
		var glr []models.Picture
		glr, imgChg = []models.Picture{}, true
		for index := range gallery {
			if _, err := url.ParseRequestURI(gallery[index]); err == nil {
				if pic := selectPic(gallery[index], cv.Gallery); pic != nil {
					glr = append(glr, *pic)
				} else {
					glr = append(glr, models.Picture{ImageURL: gallery[index]})
				}
			}
		}
		cv.Gallery = glr
		if len(cv.Gallery) > 9 { //相册最多9张
			cv.Gallery = cv.Gallery[:9]
		}
	}

	if video := c.GetString("video"); len(video) > 0 {
		cv.DesVideo, imgChg = &models.Video{VideoURL: video}, false //视频先审核再更新
		pcParam["video"] = video                                    //用户视频信息变动
		pcParam["video_check"] = models.CheckStatusUncheck          //用户视频信息变动状态
	}

	if description := c.GetString("description"); len(description) > 0 {
		param["description"], up.Description = description, description
	}

	if birth := c.GetString("birthday"); len(birth) > 0 {
		if dt, err := time.Parse("2006-01-02", birth); err == nil {
			param["birthday"], up.Birthday = dt.Unix(), dt.Unix()
		} else {
			beego.Error(err)
			dto.Message = err.Error()
			return
		}
	}

	if location := c.GetString("location"); len(location) > 0 {
		param["location"], up.Location = location, location
	}

	if price := c.GetString("price"); len(price) > 0 {
		if pr, err := strconv.ParseUint(price, 10, 64); err == nil {
			param["price"], up.Price = pr, pr
		} else {
			beego.Error("参数解析错误:"+price, err, c.Ctx.Request.UserAgent())
			dto.Message = err.Error()
			dto.Code = utils.DtoStatusParamError
			return
		}
	}

	if imgChg {
		param["cover_pic"] = cv.ToString()
	}

	trans := models.TransactionGen()
	if err := up.Update(param, trans); err != nil {
		beego.Error("更新用户资料失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "更新用户资料失败" + err.Error()
		models.TransactionRollback(trans)
		return
	}

	if err := pc.Update(pcParam, trans); err != nil {
		beego.Error("更新用户资料变动记录失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "更新用户资料变动记录失败" + err.Error()
		models.TransactionRollback(trans)
		return
	}

	upi, err := genSelfUserPorfileInfo(up, pc)
	if err != nil {
		beego.Error("获取用户资料失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "获取用户资料失败" + err.Error()
		models.TransactionRollback(trans)
		return
	}

	models.TransactionCommit(trans)

	dto.Data = upi
	dto.Sucess = true

	go checkPorn(up, cv)
}

//SendGift .
// @Title 赠送礼物
// @Description 赠送礼物
// @Param   userid     path    		string  	true        "赠送礼物的目标"
// @Param   gifid	   formData     string     	true		 "礼物id"
// @Param   count	   formData     uint     	true		 "数量"
// @Success 200 {object} utils.ResultDTO
// @router /:userid/sendgift [post]
func (c *UserController) SendGift() {
	tk, dto, uidStr := GetToken(c.Ctx), &utils.ResultDTO{}, strings.TrimSpace(c.Ctx.Input.Param(":userid"))
	defer dto.JSONResult(&c.Controller)

	giftkey, err := c.GetUint64("gifid")
	if err != nil {
		dto.Message = "必须指定礼物id"
		return
	}

	giftCount, err := c.GetUint64("count")
	if err != nil {
		beego.Error(err)
		dto.Message = "礼物数量解析失败\t" + err.Error()
		return
	}

	toID, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		beego.Error(err)
		dto.Message = "赠送礼物的目标用户ID格式错误\t" + err.Error()
		return
	}

	gft, err := (&models.Config{}).GetCommonGiftInfo()
	if err != nil {
		beego.Error(err)
		dto.Message = "获取礼物列表失败\t" + err.Error()
		return
	}

	fromUserProfile := &models.UserProfile{ID: tk.ID}
	if err = fromUserProfile.Read(); err != nil {
		beego.Error(err)
		dto.Message = "获取用户信息失败\t" + err.Error()
		return
	}

	toUserProfile := &models.UserProfile{ID: toID}

	var gift *models.Gift
	for index := range gft {
		if gft[index].ID == giftkey {
			gift = &gft[index]
		}
	}

	if gift == nil {
		beego.Error("未能找到指定的礼物" + strconv.FormatUint(giftkey, 10))
		dto.Message = "未能找到指定的礼物\t" + strconv.FormatUint(giftkey, 10)
		return
	}

	giftChg := &models.GiftChgInfo{Count: giftCount, GiftInfo: *gift}
	if err := sendGift(fromUserProfile, toUserProfile, giftChg); err != nil {
		beego.Error(err)
		dto.Message = "支付过程出现异常\t" + err.Error()
		return
	}

	dto.Sucess = true
}

//Follow .
// @Title 关注用户
// @Description 关注用户
// @Param   userid     		path    	string  	true        "用户id"
// @Success 200 {object} utils.ResultDTO
// @router /:userid/follow [post]
func (c *UserController) Follow() {
	tk, dto, uidStr := GetToken(c.Ctx), &utils.ResultDTO{}, strings.TrimSpace(c.Ctx.Input.Param(":userid"))
	defer dto.JSONResult(&c.Controller)

	toID, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		beego.Error(err)
		dto.Message = "关注用户的ID格式错误\t" + err.Error()
		return
	}

	if err := (&models.Subscribe{ID: tk.ID}).AddFollow(toID); err != nil {
		beego.Error(err)
		dto.Message = "添加关注失败\t" + err.Error()
		return
	}

	dto.Message = "关注成功"
	dto.Sucess = true
}

//UnFollow .
// @Title 取消关注用户
// @Description 取消关注用户
// @Param   userid     		path    	string  	true        "用户id"
// @Success 200 {object} utils.ResultDTO
// @router /:userid/unfollow [post]
func (c *UserController) UnFollow() {
	tk, dto, uidStr := GetToken(c.Ctx), &utils.ResultDTO{}, strings.TrimSpace(c.Ctx.Input.Param(":userid"))
	defer dto.JSONResult(&c.Controller)

	toID, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		beego.Error(err)
		dto.Message = "用户的ID格式错误\t" + err.Error()
		return
	}

	if err := (&models.Subscribe{ID: toID}).UnFollow(tk.ID); err != nil {
		beego.Error(err)
		dto.Message = "取消关注失败\t" + err.Error()
		return
	}
	dto.Message = "取消关注成功"
	dto.Sucess = true
}

//GetFollowingLst .
// @Title 获取关注列表
// @Description 获取关注列表
// @Param   length     	query    int  	false       "长度"
// @Param   skip		query    int  	false       "偏移量"
// @Success 200 {object} utils.ResultDTO
// @router /sublist [get]
func (c *UserController) GetFollowingLst() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	length, err := c.GetInt("length")
	if err != nil {
		if len(c.GetString("length")) > 0 { //length没填时默认给0,填了,但是解析错误,则返回错误
			beego.Error("参数解析错误:length\t"+err.Error(), c.Ctx.Request.UserAgent(), c.GetString("length"))
			dto.Message = "参数解析错误:length\t" + err.Error()
			dto.Code = utils.DtoStatusParamError
			return
		}
		length = 10
	}

	skip, err := c.GetInt("skip")
	if err != nil && len(c.GetString("skip")) > 0 {
		beego.Error("参数解析错误:skip\t"+err.Error(), c.Ctx.Request.UserAgent(), c.GetString("skip"))
		dto.Message = "参数解析错误:skip\t" + err.Error()
		dto.Code = utils.DtoStatusParamError
		return
	}

	up := &models.UserProfile{ID: tk.ID}
	if err := up.Read(); err != nil {
		beego.Error(err)
		dto.Message = "获取用户信息失败\t" + err.Error()
		return
	}

	sb := &models.Subscribe{ID: tk.ID}
	if err := sb.Read(nil); err != nil {
		beego.Error("获取订阅信息", err, c.Ctx.Request.UserAgent())
		dto.Message = "获取订阅信息" + err.Error()
		return
	}

	mp := sb.GetFollowing()
	var uoi []UserOperateInfo
	for k := range mp {
		if skip > 0 {
			skip--
			continue
		}

		if length > 0 {
			length--
		} else {
			break
		}

		upi := &UserPorfileInfo{UserProfile: models.UserProfile{ID: k}}
		if err := upi.Read(); err != nil {
			beego.Error("获取用户信息失败", err, c.Ctx.Request.UserAgent())
			dto.Message = "获取用户信息失败" + err.Error()
			dto.Code = utils.DtoStatusDatabaseError
			return
		}

		genUserPorfileInfoCommon(upi, upi.GetCover())
		uoi = append(uoi, UserOperateInfo{User: upi, OpTime: mp[k].FollowTime})
	}

	dto.Data = uoi
	dto.Sucess = true
}

//GetFollowedLst .
// @Title 获取粉丝列表
// @Description 获取粉丝列表
// @Param   length     	query    int  	false       "长度"
// @Param   skip		query    int  	false       "偏移量"
// @Success 200 {object} utils.ResultDTO
// @router /fanslist [get]
func (c *UserController) GetFollowedLst() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	length, err := c.GetInt("length")
	if err != nil {
		if len(c.GetString("length")) > 0 { //length没填时默认给0,填了,但是解析错误,则返回错误
			beego.Error("参数解析错误:length\t"+err.Error(), c.Ctx.Request.UserAgent(), c.GetString("length"))
			dto.Message = "参数解析错误:length\t" + err.Error()
			dto.Code = utils.DtoStatusParamError
			return
		}
		length = 10
	}

	skip, err := c.GetInt("skip")
	if err != nil && len(c.GetString("skip")) > 0 {
		beego.Error("参数解析错误:skip\t"+err.Error(), c.Ctx.Request.UserAgent(), c.GetString("skip"))
		dto.Message = "参数解析错误:skip\t" + err.Error()
		dto.Code = utils.DtoStatusParamError
		return
	}

	up := &models.UserProfile{ID: tk.ID}
	if err := up.Read(); err != nil {
		beego.Error(err)
		dto.Message = "获取用户信息失败\t" + err.Error()
		return
	}

	sb := &models.Subscribe{ID: tk.ID}
	if err := sb.Read(nil); err != nil {
		beego.Error("获取订阅信息", err, c.Ctx.Request.UserAgent())
		dto.Message = "获取订阅信息" + err.Error()
		return
	}

	mp := sb.GetFollowers()
	var uoi []UserOperateInfo
	for k := range mp {
		if skip > 0 {
			skip--
			continue
		}

		if length > 0 {
			length--
		} else {
			break
		}

		upi := &UserPorfileInfo{UserProfile: models.UserProfile{ID: k}}
		if err := upi.Read(); err != nil {
			beego.Error("获取用户信息失败", err, c.Ctx.Request.UserAgent())
			dto.Message = "获取用户信息失败" + err.Error()
			dto.Code = utils.DtoStatusDatabaseError
			return
		}

		genUserPorfileInfoCommon(upi, upi.GetCover())
		uoi = append(uoi, UserOperateInfo{User: upi, OpTime: mp[k].FollowTime})
	}

	dto.Data = uoi
	dto.Sucess = true
}

//Report .
// @Title 举报用户
// @Description 举报用户
// @Param   userid     		path    	int	  		true       "用户id"
// @Param   cate		    formData    string  	true       "举报类型"
// @Param   content     	formData    string  	true       "反馈内容"
// @Param   img		    	formData    string  	true       "图片"
// @Success 200 {object} utils.ResultDTO
// @router /:userid/report [post]
func (c *UserController) Report() {
	tk, dto, uidStr := GetToken(c.Ctx), &utils.ResultDTO{}, strings.TrimSpace(c.Ctx.Input.Param(":userid"))
	defer dto.JSONResult(&c.Controller)

	toID, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		beego.Error(err)
		dto.Message = "举报用户的ID格式错误\t" + err.Error()
		return
	}

	f := &models.FeedBack{UserID: tk.ID}
	r := &models.FeedBackReport{TgUserID: toID}
	r.Img = c.GetString("img")
	r.Content = c.GetString("content")
	r.Cate = c.GetString("cate")
	if err := f.AddReport(r); err != nil {
		beego.Error(err)
		dto.Message = "添加举报记录失败\t" + err.Error()
		return
	}

	dto.Sucess = true
}

//InviteList .
// @Title 邀请列表
// @Description 邀请列表
// @Param   length     	query    int  	false       "长度"
// @Param   skip		query    int  	false       "偏移量"
// @Success 200 {object} utils.ResultDTO
// @router /ivtlist [get]
func (c *UserController) InviteList() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	length, err := c.GetInt("length")
	if err != nil {
		if len(c.GetString("length")) > 0 { //length没填时默认给0,填了,但是解析错误,则返回错误
			beego.Error("参数解析错误:length\t"+err.Error(), c.Ctx.Request.UserAgent(), c.GetString("length"))
			dto.Message = "参数解析错误:length\t" + err.Error()
			dto.Code = utils.DtoStatusParamError
			return
		}
		length = 10
	}

	skip, err := c.GetInt("skip")
	if err != nil && len(c.GetString("skip")) > 0 {
		beego.Error("参数解析错误:skip\t"+err.Error(), c.Ctx.Request.UserAgent(), c.GetString("skip"))
		dto.Message = "参数解析错误:skip\t" + err.Error()
		dto.Code = utils.DtoStatusParamError
		return
	}

	up := &models.UserProfile{ID: tk.ID}
	lst, err := up.GetInviteList(skip, length)
	if err != nil {
		beego.Error("查询邀请列表失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "查询邀请列表失败\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	var uoi []UserOperateInfo
	for index := range lst {
		u := &models.User{ID: lst[index].ID}
		if err := u.GetRegistTime(); err != nil {
			beego.Error("获取用户信息失败", err, c.Ctx.Request.UserAgent())
			dto.Message = "获取用户信息失败" + err.Error()
			dto.Code = utils.DtoStatusDatabaseError
			return
		}

		upi := &UserPorfileInfo{UserProfile: models.UserProfile{ID: lst[index].ID}}
		if err := upi.Read(); err != nil {
			beego.Error("获取用户信息失败", err, c.Ctx.Request.UserAgent())
			dto.Message = "获取用户信息失败" + err.Error()
			dto.Code = utils.DtoStatusDatabaseError
			return
		}

		genUserPorfileInfoCommon(upi, upi.GetCover())
		uoi = append(uoi, UserOperateInfo{User: upi, OpTime: u.CreatedAt})
	}

	dto.Data = uoi
	dto.Sucess = true
}

//SetBusyStatus .
// @Title 设置状态
// @Description 设置状态
// @Param   status     	formData    int  	true       "0(勿扰),1(空闲)"
// @Success 200 {object} utils.ResultDTO
// @router /setbusy [post]
func (c *UserController) SetBusyStatus() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	status, err := c.GetInt("status")
	if err != nil {
		beego.Error("设置状态参数错误", c.Ctx.Request.UserAgent(), c.GetString("status"), err)
		dto.Message = "设置状态参数错误" + err.Error()
		dto.Code = utils.DtoStatusParamError
		return
	}

	if status == 0 {
		status = models.OnlineStatusBusy
	} else if status == 1 {
		status = models.OnlineStatusOnline
	}

	if err := (&models.UserProfile{ID: tk.ID}).UpdateOnlineStatus(status); err != nil {
		beego.Error("更新在线状态失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "更新在线状态失败\t" + err.Error()
		return
	}

	dto.Sucess = true
	dto.Message = "设置成功"
}

//BindPayAcct .
// @Title 绑定支付宝账号
// @Description 绑定支付宝账号
// @Param   acct     	formData    string  	true       "支付宝账号"
// @Param   name     	formData    string  	true       "姓名"
// @Success 200 {object} utils.ResultDTO
// @router /bindacct [post]
func (c *UserController) BindPayAcct() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	acct := c.GetString("acct")
	if len(acct) == 0 {
		dto.Message = "必须填写支付宝账号"
		return
	}

	name := c.GetString("name")
	if len(name) == 0 {
		dto.Message = "必须填写姓名"
		return
	}

	if err := (&models.UserProfile{ID: tk.ID}).UpdatePayAcct(&models.AlipayAcctInfo{Acct: acct, Name: name}); err != nil {
		beego.Error("更新提现账号失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "更新提现账号失败\t" + err.Error()
		return
	}

	dto.Message = "绑定成功"
	dto.Sucess = true
}

//AnchorApply .
// @Title 申请主播
// @Description 申请主播
// @Param   pic     	formData    string  	true       "形象照"
// @Param   video		formData    string  	true       "视频"
// @Success 200 {object} utils.ResultDTO
// @router /acapply [post]
func (c *UserController) AnchorApply() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	pic := c.GetString("pic")
	video := c.GetString("video")

	if len(pic) == 0 || len(video) == 0 {
		beego.Error("申请主播，参数错误", pic, video, c.Ctx.Request.UserAgent(), c.Ctx.Request.URL)
		dto.Message = "认证主播需要上传形象照和视频"
		dto.Code = utils.DtoStatusParamError
		return
	}

	trans := models.TransactionGen()
	if err := (&models.UserProfile{ID: tk.ID}).Update(map[string]interface{}{"anchor_auth_status": models.AnchorAuthStatusChecking}, trans); err != nil {
		models.TransactionRollback(trans)
		beego.Error("更新主播申请状态失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "更新主播申请状态失败\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	pc := &models.ProfileChg{ID: tk.ID}
	field := map[string]interface{}{"cover_pic": pic}
	field["video"] = video
	field["cover_pic_check"] = models.CheckStatusUncheck
	field["video_check"] = models.CheckStatusUncheck

	if err := pc.Update(field, trans); err != nil {
		models.TransactionRollback(trans)
		beego.Error("更新资料变动失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "更新资料变动失败\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	models.TransactionCommit(trans)
	dto.Message = "申请成功，请等待审核"
	dto.Sucess = true
}

//GuestList .
// @Title 获取访客记录
// @Description 获取访客记录
// @Param   userid     	path     int	true        "用户id,填me表示当前用户"
// @Param   length     	query    int  	false       "长度"
// @Param   skip		query    int  	false       "偏移量"
// @Success 200 {object} utils.ResultDTO
// @router /:userid/guests [get]
func (c *UserController) GuestList() {
	dto, tk, uidStr := utils.ResultDTO{}, GetToken(c.Ctx), c.Ctx.Input.Param(":userid")
	defer dto.JSONResult(&c.Controller)

	var uid uint64
	if uidStr == "me" {
		uid = tk.ID
	} else {
		var err error
		if uid, err = strconv.ParseUint(uidStr, 10, 64); err != nil {
			beego.Error("参数解析错误:userid\t"+err.Error(), c.Ctx.Request.UserAgent(), uidStr)
			dto.Message = "参数解析错误:userid\t" + err.Error()
			dto.Code = utils.DtoStatusParamError
			return
		}
	}

	length, err := c.GetInt("length")
	if err != nil {
		if len(c.GetString("length")) > 0 { //length没填时默认给0,填了,但是解析错误,则返回错误
			beego.Error("参数解析错误:length\t"+err.Error(), c.Ctx.Request.UserAgent(), c.GetString("length"))
			dto.Message = "参数解析错误:length\t" + err.Error()
			dto.Code = utils.DtoStatusParamError
			return
		}
		length = 10
	}

	skip, err := c.GetInt("skip")
	if err != nil && len(c.GetString("skip")) > 0 {
		beego.Error("参数解析错误:skip\t"+err.Error(), c.Ctx.Request.UserAgent(), c.GetString("skip"))
		dto.Message = "参数解析错误:skip\t" + err.Error()
		dto.Code = utils.DtoStatusParamError
		return
	}

	lst, err := (&models.Guest{}).GetGuestList(uid, length, skip)
	if err != nil {
		beego.Error("获取访客列表失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "获取访客列表失败\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	var ups []UserOperateInfo
	for index := range lst {
		flu := &UserPorfileInfo{UserProfile: models.UserProfile{ID: lst[index].GuestID}}
		flu.Read()
		genUserPorfileInfoCommon(flu, flu.GetCover())
		ups = append(ups, UserOperateInfo{User: flu, OpTime: lst[index].Time})
	}

	dto.Data = ups
	dto.Message = "查询成功"
	dto.Sucess = true
	return
}

//ReduceAmount .
// @Title 扣款
// @Description 扣款
// @Param   type     	formData    int  		true       "扣款类型(0:消息,1:视频)"
// @Param   target     	formData    int  		true       "目标用户ID"
// @Param   url     	formData    string  	false      "视频地址,扣款类型为1时必填"
// @Success 200 {object} utils.ResultDTO
// @router /cutamount [post]
func (c *UserController) ReduceAmount() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	var amount int
	dType, err := c.GetInt("type")
	if err != nil {
		beego.Error(err, c.Ctx.Request.UserAgent())
		dto.Message = "参数错误\t" + err.Error()
		dto.Code = utils.DtoStatusParamError
		return
	}

	tid, err := c.GetUint64("target")
	if err != nil {
		beego.Error(err, c.Ctx.Request.UserAgent())
		dto.Message = "参数错误\t" + err.Error()
		dto.Code = utils.DtoStatusParamError
		return
	}

	uri := c.GetString("url")
	if dType == 1 && len(uri) == 0 {
		dto.Message = "扣款类型为视频付费时，必须指定视频地址"
		dto.Code = utils.DtoStatusParamError
		return
	}

	if dType == 0 {
		amount = 10
	} else if dType == 1 {
		payed, err := (&models.BalanceChg{UserID: tid}).IsVideoPayed(uri, tid)
		if err != nil {
			beego.Error("获取消费记录失败", err, c.Ctx.Request.UserAgent())
			dto.Message = "获取消费记录失败"
			dto.Code = utils.DtoStatusDatabaseError
			return
		}

		if payed { //已付过费
			dto.Message = "该视频已付费,无需再付费"
			dto.Sucess = true
			return
		}

		amount = 20
	}

	up := &models.UserProfile{ID: tk.ID}
	if err := up.Read(); err != nil {
		beego.Error(err, c.Ctx.Request.UserAgent())
		dto.Message = "获取目标用户信息失败\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	mc := &models.MessageOrVideoChgInfo{TargetID: tid, URL: uri}
	mcstr, _ := utils.JSONMarshalToString(mc)

	chg := &models.BalanceChg{UserID: tk.ID, Amount: -amount, ChgInfo: mcstr}
	if dType == 0 {
		chg.ChgType = models.BalanceChgTypeMessage
	} else if dType == 1 {
		chg.ChgType = models.BalanceChgTypeVideoView
	} else {
		beego.Error("未知扣款类型", c.Ctx.Request.UserAgent(), dType)
		dto.Message = "未知扣款类型\t" + strconv.Itoa(dType)
		dto.Code = utils.DtoStatusParamError
		return
	}

	trans := models.TransactionGen() //开始事务
	if err := up.DeFund(uint64(amount), trans); err != nil {
		models.TransactionRollback(trans)
		beego.Error(err, c.Ctx.Request.UserAgent())
		dto.Message = "扣款失败\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	if err := chg.Add(trans); err != nil {
		models.TransactionRollback(trans)
		beego.Error(err, c.Ctx.Request.UserAgent())
		dto.Message = "添加余额变动失败\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	models.TransactionCommit(trans)

	dto.Message = "扣款成功"
	dto.Sucess = true
	dto.Data = chg
}

//Withdraw .
// @Title 提现
// @Description 提现
// @Param   amount     	formData    int  	true       "提现金额"
// @Success 200 {object} utils.ResultDTO
// @router /withdraw [post]
func (c *UserController) Withdraw() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	amount, err := c.GetUint64("amount")
	if err != nil {
		beego.Error(err, c.Ctx.Request.UserAgent())
		dto.Message = "参数错误\t" + err.Error()
		dto.Code = utils.DtoStatusParamError
		return
	}

	up := &models.UserProfile{ID: tk.ID}
	if err := up.Read(); err != nil {
		beego.Error(err, c.Ctx.Request.UserAgent())
		dto.Message = "获取目标用户信息失败\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	if up.Income < amount {
		beego.Error("账户余额不足以提现:"+strconv.FormatUint(up.Income, 10)+strconv.FormatUint(amount, 10), c.Ctx.Request.UserAgent())
		dto.Message = "账户余额不足以提现"
		dto.Code = utils.DtoStatusParamError
		return
	}

	trans := models.TransactionGen()
	if err := up.AddIncome(-int(amount), trans); err != nil {
		beego.Error("申请提现扣款失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "申请提现扣款失败" + err.Error()
		models.TransactionRollback(trans)
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	wd := &models.Withdraw{UserID: tk.ID, Amount: amount, CreateAt: time.Now().Unix()}
	if err := wd.Add(trans); err != nil {
		beego.Error("生成提现申请失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "生成提现申请失败" + err.Error()
		models.TransactionRollback(trans)
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	wdStr, err := utils.JSONMarshalToString(wd)
	if err != nil {
		beego.Error("解析变动信息失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "解析变动信息失败" + err.Error()
		models.TransactionRollback(trans)
		return
	}
	bc := &models.BalanceChg{UserID: tk.ID, ChgType: models.BalanceChgTypeWithDraw, ChgInfo: wdStr, Amount: -int(amount), Time: time.Now().Unix()}
	if err := bc.Add(trans); err != nil {
		beego.Error("生成变动记录失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "生成变动记录失败" + err.Error()
		models.TransactionRollback(trans)
		return
	}

	models.TransactionCommit(trans)
	dto.Sucess = true
	dto.Message = "提现申请已提交"
}

//WithdrawHis .
// @Title 获取提现记录
// @Description 获取提现记录
// @Param   length     	query    int  	false       "长度"
// @Param   skip		query    int  	false       "偏移量"
// @Success 200 {object} utils.ResultDTO
// @router /withdrawhis [get]
func (c *UserController) WithdrawHis() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	length, err := c.GetInt("length")
	if err != nil {
		if len(c.GetString("length")) > 0 { //length没填时默认给0,填了,但是解析错误,则返回错误
			beego.Error("参数解析错误:length\t"+err.Error(), c.Ctx.Request.UserAgent(), c.GetString("length"))
			dto.Message = "参数解析错误:length\t" + err.Error()
			dto.Code = utils.DtoStatusParamError
			return
		}
		length = 10
	}

	skip, err := c.GetInt("skip")
	if err != nil && len(c.GetString("skip")) > 0 {
		beego.Error("参数解析错误:skip\t"+err.Error(), c.Ctx.Request.UserAgent(), c.GetString("skip"))
		dto.Message = "参数解析错误:skip\t" + err.Error()
		dto.Code = utils.DtoStatusParamError
		return
	}

	wds, err := (&models.Withdraw{UserID: tk.ID}).List(skip, length)
	if err != nil {
		beego.Error("查询提现记录失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "查询提现记录失败\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	dto.Data = wds
	dto.Sucess = true
	dto.Message = "查询成功"
}

//GiftHistory .
// @Title 获取用户收到的所有礼物
// @Description 获取用户收到的所有礼物
// @Param   userid     path    string  true        "用户id,填me表示获取当前账号的用户信息"
// @Success 200 {object} models.UserProfile
// @router /:userid/gifthis [get]
func (c *UserController) GiftHistory() {
	dto := &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	var uid uint64
	var err error
	uidStr := strings.TrimSpace(c.Ctx.Input.Param(":userid"))
	if uidStr == "me" {
		tk := GetToken(c.Ctx)
		uid = tk.ID
	} else {
		uid, err = strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			beego.Error(err)
			dto.Message = err.Error()
			return
		}
	}

	ue := &models.UserExtra{ID: uid}
	lst, err := ue.GetGiftHis()
	if err != nil {
		beego.Error("获取礼物历史记录失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "获取礼物历史记录失败\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	dto.Data = lst
	dto.Sucess = true
	dto.Message = "查询成功"
}

//赠送礼物,流程包括 源用户扣款，目标用户增加余额，邀请人分成，以及分别添加余额变动记录,过程中任何一部出错，事务回滚并返回失败
//赠送礼物不参与分成
func sendGift(from, to *models.UserProfile, gift *models.GiftChgInfo) error {
	u := &models.User{ID: to.ID}
	if err := u.Read(); err != nil {
		return errors.New("获取目标用户信息失败,id:" + strconv.FormatUint(to.ID, 10) + "\t" + err.Error())
	}

	chgInfo, _ := utils.JSONMarshalToString(gift)

	amount := gift.GiftInfo.Price * gift.Count         //消费金额
	income, inviteIncome, err := computeIncome(amount) //收益金额,分成金额
	if err != nil {
		return err
	}

	iu, iuchg := genInvitationIncome(to.ID, u.InvitedBy, inviteIncome, chgInfo)

	trans := models.TransactionGen() //开始事务
	if err := from.AllocateFund(to, iu, amount, uint64(income), uint64(inviteIncome), trans); err != nil {
		models.TransactionRollback(trans)
		return err
	}

	if err := (&models.UserExtra{ID: to.ID}).AddGiftCount(gift.GiftInfo, gift.Count, trans); err != nil {
		models.TransactionRollback(trans)
		return err
	}

	//增加历史收益，邀请人历史收益
	if err := (&models.UserExtra{ID: to.ID}).AddIncomeHis(uint64(income), trans); err != nil {
		models.TransactionRollback(trans)
		return errors.New("增加赠礼目标用户历史收益失败\t" + err.Error())
	}

	if u.InvitedBy != 0 { //如果有邀请人，增加邀请人历史收益
		if err := (&models.UserExtra{ID: u.InvitedBy}).AddInviteIncomeHis(uint64(inviteIncome), trans); err != nil {
			models.TransactionRollback(trans)
			return errors.New("增加赠礼目标用户邀请人历史收益失败\t" + err.Error())
		}
	}

	fuChg := &models.BalanceChg{UserID: from.ID, FromUserID: to.ID, ChgType: models.BalanceChgTypeGift, Amount: -int(amount)} //源用户扣款变动
	fuChg.ChgInfo = chgInfo

	tuchg := &models.BalanceChg{UserID: to.ID, FromUserID: from.ID, ChgType: models.BalanceChgTypeGift, Amount: income} //目标用户余额增加 变动
	tuchg.ChgInfo = chgInfo

	if err := fuChg.AddChg(trans, fuChg, tuchg, iuchg); err != nil {
		models.TransactionRollback(trans)
		return err
	}

	models.TransactionCommit(trans) //提交事务
	return nil
}

//视频聊天结算,增加主播收入和邀请人收入
func videoDone(from, to *models.UserProfile, video *models.VideoChgInfo, amount uint64) error {
	u := &models.User{ID: to.ID}
	if err := u.Read(); err != nil {
		return errors.New("获取目标用户信息失败,id:" + strconv.FormatUint(to.ID, 10) + "\t" + err.Error())
	}

	chgInfo, _ := utils.JSONMarshalToString(video)

	income, inviteIncome, err := computeIncome(amount) //收益金额,分成金额
	if err != nil {
		return err
	}

	iu, iuchg := genInvitationIncome(to.ID, u.InvitedBy, inviteIncome, chgInfo)

	trans := models.TransactionGen() //开始事务
	//视频结算不扣除用户费用(已在websocket连接时每分钟扣费时结清)，只增加主播收益和主播邀请人收益

	if err := (&models.UserProfile{ID: to.ID}).AddIncome(income, trans); err != nil {
		models.TransactionRollback(trans)
		return err
	}

	if iu != nil {
		if err := iu.AddBalance(inviteIncome, trans); err != nil { //邀请人收益
			models.TransactionRollback(trans)
			return err
		}
	}

	//增加历史收益，邀请人历史收益
	if err := (&models.UserExtra{ID: to.ID}).AddIncomeHis(uint64(income), trans); err != nil {
		models.TransactionRollback(trans)
		return errors.New("增加赠礼目标用户历史收益失败\t" + err.Error())
	}

	if u.InvitedBy != 0 { //如果有邀请人，增加邀请人历史收益
		if err := (&models.UserExtra{ID: u.InvitedBy}).AddInviteIncomeHis(uint64(inviteIncome), trans); err != nil {
			models.TransactionRollback(trans)
			return errors.New("增加赠礼目标用户邀请人历史收益失败\t" + err.Error())
		}
	}

	fuChg := &models.BalanceChg{UserID: from.ID, FromUserID: to.ID, ChgType: models.BalanceChgTypeVideo, Amount: -int(amount)} //源用户扣款变动
	fuChg.ChgInfo = chgInfo

	tuchg := &models.BalanceChg{UserID: to.ID, FromUserID: from.ID, ChgType: models.BalanceChgTypeVideo, Amount: income} //目标用户余额增加 变动
	tuchg.ChgInfo = chgInfo

	if err := fuChg.AddChg(trans, fuChg, tuchg, iuchg); err != nil {
		models.TransactionRollback(trans)
		return err
	}

	models.TransactionCommit(trans) //提交事务
	return nil
}

//videoAllocateFund 视频聊天定时扣费 只扣用户的钱，主播和主播邀请人的收益留到视频结束时统一增加（避免计算收益率时产生的精度误差）
func videoAllocateFund(from, to *models.UserProfile, price uint64) error {
	trans := models.TransactionGen() //开始事务

	beego.Info("视频扣费,用户ID:", from.ID, "金额:", price)
	if err := from.Read(); err != nil {
		return err
	}

	if err := from.DeFund(price, trans); err != nil {
		models.TransactionRollback(trans)
		return err
	}

	models.TransactionCommit(trans) //提交事务
	return nil
}

//computeIncome 计算收益
func computeIncome(amount uint64) (income, inviteIncome int, err error) {
	rate, err := (&models.Config{}).GetIncomeRate()
	if err != nil {
		err = errors.New("获取收益分成率失败" + err.Error())
		return
	}

	income = int(float64(amount) * (1 - rate.IncomeFee))            //收益金额
	inviteIncome = int(float64(income) * (rate.InviteIncomegeRate)) //分成金额
	return
}

//genInvitationIncome 生成邀请人收益变动
func genInvitationIncome(uid, invitedByid uint64, inviteIncome int, chgInfo string) (iu *models.UserProfile, iuchg *models.BalanceChg) {
	if invitedByid != 0 { //邀请人分成
		iuchg = &models.BalanceChg{UserID: invitedByid, FromUserID: uid, ChgType: models.BalanceChgTypeInvitationIncome, Amount: inviteIncome}
		iuchg.ChgInfo = chgInfo
		iu = &models.UserProfile{ID: invitedByid}
	}
	return
}

func selectPic(imgURL string, arr []models.Picture) *models.Picture {
	for index := range arr {
		if arr[index].ImageURL == imgURL {
			return &arr[index]
		}
	}
	return nil
}

func checkPorn(up *models.UserProfile, cover *models.UserCoverInfo) {

}
