package models

import (
	"MoShow/utils"
	"errors"
	"strconv"
	"time"

	"github.com/astaxie/beego"

	"github.com/jinzhu/gorm"
)

const (
	//GenderDefault .
	GenderDefault = iota
	//GenderMan .
	GenderMan
	//GenderWoman .
	GenderWoman
)

const (
	//UserTypeNormal 普通用户
	UserTypeNormal = iota
	//UserTypeAnchor 主播
	UserTypeAnchor
	//UserTypeFaker 马甲号
	UserTypeFaker
)

const (
	//UserStatusNormal 正常
	UserStatusNormal = iota
	//UserStatusHot 推荐
	UserStatusHot
	//UserStatusBlock 屏蔽
	UserStatusBlock
)

const (
	//OnlineStatusOffline 离线
	OnlineStatusOffline = iota
	//OnlineStatusOnline 在线
	OnlineStatusOnline
	//OnlineStatusChating 正在聊天
	OnlineStatusChating
	//OnlineStatusBusy 勿扰
	OnlineStatusBusy
)

const (
	//AnchorAuthStatusUnAuth 未认证
	AnchorAuthStatusUnAuth = iota
	//AnchorAuthStatusChecking 审核中
	AnchorAuthStatusChecking
	//AnchorAuthStatusDone 审核完成
	AnchorAuthStatusDone
)

//UserProfile .
type UserProfile struct {
	ID               uint64  `json:"user_id" gorm:"column:id;primary_key" description:"用户ID"`
	Alias            string  `json:"alias" gorm:"column:alias" description:"昵称"`
	Gender           int     `json:"gender" gorm:"column:gender" description:"性别"`
	CoverPic         string  `json:"-" gorm:"column:cover" description:"形象展示,包括头像,相册,视频"`
	Description      string  `json:"description" gorm:"column:description" description:"签名"`
	Birthday         int64   `json:"birthday" gorm:"column:birthday" description:"生日"`
	Location         string  `json:"location" gorm:"column:location" description:"地区"`
	Balance          uint64  `json:"balance,omitempty" gorm:"column:balance" description:"余额"`
	Income           uint64  `json:"income,omitempty" gorm:"column:income" description:"收益"`
	Price            uint64  `json:"price,omitempty" gorm:"column:price" description:"视频价格/分"`
	UserType         int     `json:"user_type" gorm:"column:user_type" description:"用户类型"`
	ImToken          string  `json:"-" gorm:"column:im_token" description:"网易云信token"`
	UserStatus       int     `json:"-" gorm:"column:user_status" description:"用户状态"`
	OnlineStatus     int     `json:"online_status" gorm:"column:online_status" description:"在线状态"`
	AnchorAuthStatus int     `json:"anchor_auth_status" gorm:"column:anchor_auth_status" description:"主播认证状态"`
	DialAccept       uint64  `json:"-" gorm:"column:dial_accept" description:"视频接通数"`
	DialDeny         uint64  `json:"-" gorm:"column:dial_deny" description:"视频拒接数"`
	DialDuration     uint64  `json:"dial_duration" gorm:"column:dial_duration" description:"通话时间"`
	Duration         uint64  `json:"recent_duration" gorm:"column:recent_duration" description:"最近通话时间"`
	AlipayAcct       string  `json:"-" gorm:"column:alipay_acct" description:"支付宝账号"`
	SpecialRate      float64 `json:"-" gorm:"column:special_rate" description:"特殊分成"`
	UpdateAt         int64   `json:"update_at" gorm:"column:update_at" description:"更新时间"`
}

//UserCoverInfo .
type UserCoverInfo struct {
	CoverPicture *Picture  `json:"cover_pic_info"`
	DesVideo     *Video    `json:"video_info"`
	Gallery      []Picture `json:"gallery"`
}

//Picture .
type Picture struct {
	ImageURL   string `json:"image_url"`
	CloudCheck bool   `json:"cloud_porn_check"`
}

//Video .
type Video struct {
	VideoURL string `json:"video_url"`
}

//AlipayAcctInfo .
type AlipayAcctInfo struct {
	Acct string `json:"acct"`
	Name string `json:"name"`
}

//TableName .
func (UserProfile) TableName() string {
	return "user_profile"
}

//ToString .
func (u *UserCoverInfo) ToString() string {
	str, err := utils.JSONMarshalToString(u)
	if err != nil {
		beego.Error("用户头像，相册解析出错", err)
		return "{}"
	}
	return str
}

//Add .
func (u *UserProfile) Add(trans *gorm.DB) error {
	if u.ID == 0 {
		return errors.New("必须指定用户ID")
	}

	if trans == nil {
		trans = db
	}

	if len(u.AlipayAcct) == 0 {
		u.AlipayAcct = "{}"
	}

	if len(u.CoverPic) == 0 {
		u.CoverPic = "{}"
	}

	return trans.Create(u).Error
}

//Read .
func (u *UserProfile) Read() error {
	if u.ID == 0 {
		return errors.New("必须指定user_profile的id")
	}

	return db.Find(u).Error
}

//ReadWallet .
func (u *UserProfile) ReadWallet() error {
	return db.Select("income, balance").Find(u).Error
}

//Update .
func (u *UserProfile) Update(fields map[string]interface{}, trans *gorm.DB) error {
	fields["update_at"] = time.Now().Unix()

	if len(fields) == 1 {
		return nil
	}

	if trans == nil {
		trans = db
	}
	return trans.Model(u).Updates(fields).Error
}

//UpdateRecentDialTime .
func (u *UserProfile) UpdateRecentDialTime() error {
	return db.Model(u).Update("recent_duration", gorm.Expr("COALESCE((SELECT sum(duration) FROM MoShow.dial where (to_user_id = ? or from_user_id = ?) and create_at >= ? GROUP BY to_user_id or from_user_id),0)", u.ID, u.ID, time.Now().Add(-3*24*time.Hour).Unix())).Error
}

//UpdateCover .
func (u *UserProfile) UpdateCover(cover *UserCoverInfo) error {
	cs, err := utils.JSONMarshalToString(cover)
	if err != nil {
		return err
	}
	return db.Model(u).Updates(map[string]interface{}{"cover": cs}).Error
}

//UpdatePayAcct .
func (u *UserProfile) UpdatePayAcct(acct *AlipayAcctInfo) error {
	acctStr, _ := utils.JSONMarshalToString(acct)
	return db.Model(u).Update("alipay_acct", acctStr).Error
}

//UpdateOnlineStatus .
func (u *UserProfile) UpdateOnlineStatus(status int) error {
	return db.Model(u).Update("online_status", status).Error
}

//SetFaker .
func (u *UserProfile) SetFaker() error {
	return db.Model(u).Updates(map[string]interface{}{"user_type": UserTypeFaker, "price": 0}).Error
}

//ResetOnlineStatus .
func (u *UserProfile) ResetOnlineStatus() error {
	return db.Model(u).Where("online_status  = ?", OnlineStatusChating).Update("online_status", OnlineStatusOnline).Error
}

//GetCover .
func (u *UserProfile) GetCover() *UserCoverInfo {
	if len(u.CoverPic) > 0 {
		ucp := &UserCoverInfo{}
		if err := utils.JSONUnMarshal(u.CoverPic, ucp); err != nil {
			beego.Error(err)
			return nil
		}
		return ucp
	}
	return nil
}

//GetInviteList 获取我邀请的用户列表
func (u *UserProfile) GetInviteList(skip, limit int) (ul []UserProfile, err error) {
	if err = db.Joins("left join users on user_profile.id = users.id").Where("invited_by = ?", u.ID).Order("users.create_at desc").Offset(skip).Limit(limit).Find(&ul).Error; err != nil {
		return nil, err
	}

	return
}

//AddDialDuration 增加通话时间
func (u *UserProfile) AddDialDuration(duration uint64, trans *gorm.DB) error {
	if duration == 0 {
		return nil
	}

	if trans == nil {
		trans = db
	}
	return trans.Model(u).Updates(map[string]interface{}{"dial_duration": gorm.Expr("dial_duration + ?", duration), "dial_accept": gorm.Expr("dial_accept + ?", 1)}).Error
}

//AddDialReject .
func (u *UserProfile) AddDialReject(trans *gorm.DB) error {
	if trans == nil {
		trans = db
	}
	return trans.Model(u).Update("dial_deny", gorm.Expr("dial_deny + ?", 1)).Error
}

//AddBalance .
func (u *UserProfile) AddBalance(amount int, trans *gorm.DB) error {
	if amount == 0 {
		return nil
	}

	if trans == nil {
		trans = db
	}

	return trans.Model(u).Update("balance", gorm.Expr("balance + ?", amount)).Error
}

//AddIncome .
func (u *UserProfile) AddIncome(amount int, trans *gorm.DB) error {
	if amount == 0 {
		return nil
	}

	if trans == nil {
		trans = db
	}
	return trans.Model(u).Update("income", gorm.Expr("income + ?", amount)).Error
}

//DeFund 用户扣款
func (u *UserProfile) DeFund(amount uint64, trans *gorm.DB) error {
	if amount == 0 {
		return nil
	}

	if u.Balance+u.Income < amount { //检查余额
		return errors.New("用户余额不足，扣款(" + strconv.FormatUint(amount, 10) + ")失败,所有钱包余额合计:" + strconv.FormatUint(u.Balance+u.Income, 10))
	}

	if u.Balance > amount { //余额钱包金额足够扣款
		if err := u.AddBalance(-int(amount), trans); err != nil { //扣款
			return errors.New(strconv.FormatUint(u.ID, 10) + "扣款失败\t" + err.Error())
		}
		u.Balance -= amount
	} else { //余额钱包金额不足以扣款
		deIncome := amount - u.Balance
		if err := u.AddBalance(-int(u.Balance), trans); err != nil { //从余额钱包扣款
			return errors.New(strconv.FormatUint(u.ID, 10) + "扣款失败\t" + err.Error())
		}

		if err := u.AddIncome(-int(deIncome), trans); err != nil { //从收益钱包扣款
			return errors.New(strconv.FormatUint(u.ID, 10) + "扣款失败\t" + err.Error())
		}
		u.Income -= deIncome
		u.Balance = 0
	}
	return nil
}

//AllocateFund 划款
func (u *UserProfile) AllocateFund(to, invitation *UserProfile, amount, incomeAmount, inviteAmount uint64, trans *gorm.DB) error {
	if err := u.DeFund(amount, trans); err != nil {
		return err
	}

	if err := to.AddIncome(int(incomeAmount), trans); err != nil { //增加余额
		return errors.New("接受人增加余额失败\t" + err.Error())
	}

	if invitation != nil {
		if uint64(incomeAmount*3/10) < inviteAmount {
			return errors.New("收益分成比例异常,收入:" + strconv.FormatUint(incomeAmount, 10) + ",分成:" + strconv.FormatUint(inviteAmount, 10))
		}

		if err := invitation.AddIncome(int(inviteAmount), trans); err != nil { //增加余额
			return errors.New("邀请人增加余额失败\t" + err.Error())
		}
	}

	return nil
}
