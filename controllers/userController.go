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
	ImTk      string                `json:"im_token,omitempty"`
	CoverInfo *models.UserCoverInfo `json:"cover_info"  description:"形象展示,包括头像,相册,视频"`
	Followed  bool                  `json:"followed" description:"是否已关注"`
}

//Read .
// @Title 读取用户
// @Description 读取用户
// @Param   userid     path    string  true        "用户id,填me表示获取当前账号的用户信息"
// @Success 200 {object} models.UserProfile
// @router /:userid [get]
func (c *UserController) Read() {
	dto := utils.ResultDTO{Sucess: false}
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

		if uid != tk.ID { //自己查看自己时不增加记录
			if err := (&models.Guest{}).AddView(uid, tk.ID); err != nil { //增加访客记录
				beego.Error(err)
				dto.Message = "增加访客记录失败\t" + err.Error()
			}
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
	upi.CoverInfo = up.GetCover()
	if uid == tk.ID {
		upi.ImTk = up.ImToken
	}

	if fl := up.GetFollowers(); fl != nil {
		if _, ok := fl[tk.ID]; ok {
			upi.Followed = true
		}
	}

	dto.Data = upi
	dto.Sucess = true
}

//Update .
// @Title 更新用户
// @Description 更新用户
// @Param   alias     		formData    string  	true        "昵称"
// @Param   cover_pic	    formData    string  	false       "头像"
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
		beego.Error(err)
		dto.Message = err.Error()
		return
	}

	cv, imgChg, param := up.GetCover(), false, make(map[string]interface{})
	if cv == nil {
		cv = &models.UserCoverInfo{}
	}

	if alias := c.GetString("alias"); len(alias) > 0 {
		param["alias"], up.Alias = alias, alias
	}

	if coverPic := c.GetString("cover_pic"); len(coverPic) > 0 {
		cv.CoverPicture, imgChg = &models.Picture{ImageURL: coverPic}, true
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
		cv.DesVideo, imgChg = &models.Video{VideoURL: video}, true
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
			beego.Error(err)
			dto.Message = err.Error()
			return
		}
	}

	if imgChg {
		param["cover_pic"] = cv.ToString()
	}

	if err := up.Update(param); err != nil {
		beego.Error(err)
		dto.Message = err.Error()
		return
	}

	dto.Data = &UserPorfileInfo{UserProfile: *up, CoverInfo: cv}
	dto.Sucess = true
}

//SendGift .
// @Title 赠送礼物
// @Description 赠送礼物
// @Param   userid     path    		string  	true        "赠送礼物的目标"
// @Param   giftkey	   formData     string     	true		 "礼物id"
// @Param   count	   formData     uint     	true		 "数量"
// @Success 200 {object} utils.ResultDTO
// @router /:userid/sendgift [post]
func (c *UserController) SendGift() {
	tk, dto, uidStr := GetToken(c.Ctx), &utils.ResultDTO{}, strings.TrimSpace(c.Ctx.Input.Param(":userid"))
	defer dto.JSONResult(&c.Controller)

	giftkey := c.GetString("giftkey")
	if len(giftkey) == 0 {
		dto.Message = "必须指定礼物key"
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
	if err := toUserProfile.Read(); err != nil {
		beego.Error(err)
		dto.Message = "获取用户信息失败\t" + err.Error()
		return
	}

	gift, ok := gft[giftkey]
	if !ok {
		beego.Error("未能找到指定的礼物" + giftkey)
		dto.Message = "未能找到指定的礼物\t" + giftkey
		return
	}

	giftChg := &models.GiftChgInfo{Count: giftCount, GiftInfo: &gift}
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

	up := models.UserProfile{ID: toID}
	if err := up.AddFollow(tk.ID); err != nil {
		beego.Error(err)
		dto.Message = "添加关注失败\t" + err.Error()
		return
	}

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

	up := models.UserProfile{ID: toID}
	if err := up.UnFollow(tk.ID); err != nil {
		beego.Error(err)
		dto.Message = "取消关注失败\t" + err.Error()
		return
	}

	dto.Sucess = true
}

//GetFollowingLst .
// @Title 获取关注列表
// @Description 获取关注列表
// @Success 200 {object} utils.ResultDTO
// @router /sublist [get]
func (c *UserController) GetFollowingLst() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	up := &models.UserProfile{ID: tk.ID}
	if err := up.Read(); err != nil {
		beego.Error(err)
		dto.Message = "获取用户信息失败\t" + err.Error()
		return
	}

	mp := up.GetFollowing()
	var flst []uint64
	for k := range mp {
		flst = append(flst, k)
	}

	dto.Data = flst
	dto.Sucess = true
}

//GetFollowedLst .
// @Title 获取粉丝列表
// @Description 获取粉丝列表
// @Success 200 {object} utils.ResultDTO
// @router /fanslist [get]
func (c *UserController) GetFollowedLst() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	up := &models.UserProfile{ID: tk.ID}
	if err := up.Read(); err != nil {
		beego.Error(err)
		dto.Message = "获取用户信息失败\t" + err.Error()
		return
	}

	mp := up.GetFollowers()
	var flst []uint64
	for k := range mp {
		flst = append(flst, k)
	}

	dto.Data = flst
	dto.Sucess = true
}

//Report .
// @Title 举报用户
// @Description 举报用户
// @Param   userid     		path    	int	  		true        "用户id"
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
// @Success 200 {object} utils.ResultDTO
// @router /ivtlist [get]
func (c *UserController) InviteList() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	up := &models.UserProfile{ID: tk.ID}
	lst, err := up.GetInviteList()
	if err != nil {
		beego.Error("查询邀请列表失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "查询邀请列表失败\t" + err.Error()
		dto.Code = utils.DtoStatusDatabaseError
		return
	}

	dto.Data = lst
	dto.Sucess = true
}

//ReduceAmount .
// @Title 扣款
// @Description 扣款
// @Param   type     	formData    int  	true       "扣款类型(0:消息,1:视频)"
// @Param   amount     	formData    int  	true       "扣款金额"
// @Success 200 {object} utils.ResultDTO
// @router /cutamount [post]
func (c *UserController) ReduceAmount() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	amount, err := c.GetUint64("amount")
	if err != nil {
		beego.Error(err, c.Ctx.Request.UserAgent())
		dto.Message = "参数错误\t" + err.Error()
		dto.Code = utils.DtoStatusParamError
		return
	}

	dType, err := c.GetInt("type")
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

	chg := &models.BalanceChg{UserID: tk.ID, Amount: -int(amount), ChgInfo: "{}"}
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
	if err := up.DeFund(amount, trans); err != nil {
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

func videoAllocateFund(from, to *models.UserProfile, price uint64) error {
	income, inviteIncome, err := computeIncome(price) //收益金额,分成金额
	if err != nil {
		return err
	}

	trans := models.TransactionGen() //开始事务

	if err := from.AllocateFund(to, nil, price, uint64(income), uint64(inviteIncome), trans); err != nil {
		models.TransactionRollback(trans)
		return err
	}

	models.TransactionCommit(trans) //提交事务
	return nil
}

//视频聊天结算
func videoDone(from, to *models.UserProfile, video *models.VideoChgInfo) error {
	u := &models.User{ID: to.ID}
	if err := u.Read(); err != nil {
		return errors.New("获取目标用户信息失败,id:" + strconv.FormatUint(to.ID, 10) + "\t" + err.Error())
	}

	chgInfo, _ := utils.JSONMarshalToString(video)

	amount := video.Price * video.TimeLong             //消费金额
	income, inviteIncome, err := computeIncome(amount) //收益金额,分成金额
	if err != nil {
		return err
	}

	iu, iuchg := genInvitationIncome(to.ID, u.InvitedBy, inviteIncome, chgInfo)

	trans := models.TransactionGen() //开始事务
	//视频结算只生成变动，金额已在websocket中扣除

	if err := iu.AddBalance(inviteIncome, trans); err != nil { //邀请人收益
		models.TransactionRollback(trans)
		return err
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

//计算收益
func computeIncome(amount uint64) (income, inviteIncome int, err error) {
	rate, err := (&models.Config{}).GetIncomeRate()
	if err != nil {
		err = errors.New("获取收益分成率失败" + err.Error())
		return
	}

	income = int(float64(amount) * (rate.IncomeFee))                //收益金额
	inviteIncome = int(float64(income) * (rate.InviteIncomegeRate)) //分成金额
	return
}

//生成邀请人收益变动
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
