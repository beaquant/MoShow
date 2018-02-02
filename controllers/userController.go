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
	CoverInfo *models.UserCoverInfo `json:"cover_info"  description:"形象展示,包括头像,相册,视频"`
}

//Create .
// @Title 创建用户信息
// @Description 创建用户信息
// @Param   alias     		formData    string  	true        "昵称"
// @Param   gender    		formData    int     	true        "性别，男(1),女(0)"
// @Param   cover_pic	    formData    string  	false       "头像"
// @Param   gallery		    formData    string  	false       "相册"
// @Param   description     formData    string  	false       "签名"
// @Param   birthday     	formData    time.Time  	false       "生日,格式:2006-01-02"
// @Param   location     	formData    string  	false       "地区"
// @Param   price	     	formData    uint	  	false       "价格"
// @Success 200 {object} 	utils.ResultDTO
// @router /create [put]
func (c *UserController) Create() {
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	up := &models.UserProfile{ID: tk.ID}
	if alias := c.GetString("alias"); len(alias) > 0 {
		up.Alias = alias
	}

	if gender := c.GetString("gender"); len(gender) > 0 {
		if gd, err := strconv.Atoi(gender); err == nil && gd == 0 || gd == 1 {
			up.Gender = gd
		}
	}

	if description := c.GetString("description"); len(description) > 0 {
		up.Description = description
	}

	if birth := c.GetString("birthday"); len(birth) > 0 {
		if dt, err := time.Parse("2006-01-02", birth); err == nil {
			up.Birthday = dt.Unix()
		} else {
			beego.Error(err)
			dto.Message = err.Error()
			return
		}
	} else {
		up.Birthday = time.Date(1993, 1, 1, 0, 0, 0, 0, nil).Unix()
	}

	if location := c.GetString("location"); len(location) > 0 {
		up.Location = location
	}

	if price := c.GetString("price"); len(price) > 0 {
		if pr, err := strconv.ParseUint(price, 10, 64); err == nil {
			up.Price = pr
		} else {
			beego.Error(err)
			dto.Message = err.Error()
			return
		}
	}

	uci := &models.UserCoverInfo{}
	if coverPic := c.GetString("cover_pic"); len(coverPic) > 0 {
		uci.CoverPicture = &models.Picture{ImageURL: coverPic}
	}

	if gallery := c.GetStrings("gallery"); gallery != nil && len(gallery) > 0 {
		for index := range gallery {
			if _, err := url.ParseRequestURI(gallery[index]); err == nil {
				uci.Gallery = append(uci.Gallery, models.Picture{ImageURL: gallery[index]})
			}
		}
		if len(uci.Gallery) > 9 {
			uci.Gallery = uci.Gallery[:9]
		}
	}

	up.CoverPic = uci.ToString()
	err := up.Add()
	if err != nil {
		beego.Error(err)
		dto.Message = err.Error()
		return
	}

	go checkPorn(up, uci) //鉴黄

	dto.Sucess = true
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

//Report .
// @Title 举报用户
// @Description 举报用户
// @Param   userid     		path    	string  	true        "用户id"
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

//赠送礼物,流程包括 源用户扣款，目标用户增加余额，邀请人分成，以及分别添加余额变动记录,过程中任何一部出错，事务回滚并返回失败
//赠送礼物不参与分成
func sendGift(from, to *models.UserProfile, gift *models.GiftChgInfo) error {
	amount := uint64(gift.GiftInfo.Price) * gift.Count
	if from.Balance < amount { //检查余额
		return errors.New("用户余额不足，送礼失败,余额:" + strconv.FormatUint(from.Balance, 64))
	}

	trans := models.TransactionGen() //开始事务
	if err := from.AllocateFund(to, nil, amount, 0, trans); err != nil {
		models.TransactionRollback(trans)
		return err
	}

	fuChg := &models.BalanceChg{UserID: from.ID, FromUserID: to.ID, ChgType: models.BalanceChgTypeSendGift, Amount: -int(amount)} //源用户扣款变动
	fuChg.ChgInfo, _ = utils.JSONMarshalToString(gift)

	tuchg := &models.BalanceChg{UserID: to.ID, FromUserID: from.ID, ChgType: models.BalanceChgTypeReceiveGift, Amount: int(amount)} //目标用户余额增加 变动
	tuchg.ChgInfo, _ = utils.JSONMarshalToString(gift)

	if err := fuChg.AddChg(trans, fuChg, tuchg); err != nil {
		models.TransactionRollback(trans)
		return err
	}

	models.TransactionCommit(trans) //提交事务
	return nil
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
