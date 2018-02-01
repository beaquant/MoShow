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
	//GenderWoman .
	GenderWoman = iota
	//GenderMan .
	GenderMan
)

const (
	//UserTypeNormal 普通用户
	UserTypeNormal = iota
	//UserTypeAnchor 主播
	UserTypeAnchor
)

const (
	//UserStatusNormal 正常
	UserStatusNormal = iota
	//UserStatusHot 推荐
	UserStatusHot
	//UserStatusBlock 屏蔽
	UserStatusBlock
)

//UserProfile .
type UserProfile struct {
	ID          uint64 `json:"user_id" gorm:"column:id;primary_key" description:"用户ID"`
	Alias       string `json:"alias" gorm:"column:alias" description:"昵称"`
	Gender      int    `json:"gender" gorm:"column:gender" description:"性别"`
	CoverPic    string `json:"-" gorm:"column:cover" description:"形象展示,包括头像,相册,视频"`
	Description string `json:"description" gorm:"column:description" description:"签名"`
	Birthday    int64  `json:"birthday" gorm:"column:birthday" description:"生日"`
	Location    string `json:"location" gorm:"column:location" description:"地区"`
	Balance     uint64 `json:"balance" gorm:"column:balance" description:"余额"`
	Price       uint64 `json:"price" gorm:"column:price" description:"视频价格/分"`
	UserType    int    `json:"user_type" gorm:"column:user_type" description:"用户类型"`
	UserStatus  int    `json:"-" gorm:"column:user_status" description:"用户状态"`
	UpdateAt    int64  `json:"update_at" gorm:"column:update_at" description:"更新时间"`
}

//UserCoverInfo .
type UserCoverInfo struct {
	CoverPicture *Picture  `json:"cover_pic_info"`
	DesVideo     *Video    `json:"video_info"`
	Gallery      []Picture `json:"gallery"`
}

//Picture .
type Picture struct {
	ImageURL string `json:"image_url"`
	Disable  bool   `json:"disable"`
	Checked  bool   `json:"checked"`
}

//Video .
type Video struct {
	VideoURL string `json:"video_url"`
	Disable  bool   `json:"disable"`
	Checked  bool   `json:"checked"`
}

//TableName .
func (UserProfile) TableName() string {
	return "user_profile"
}

//ToString .
func (u *UserCoverInfo) ToString() string {
	str, err := utils.JSONMarshalToString(u)
	if err != nil {
		beego.Error(err)
		return ""
	}
	return str
}

//Add .
func (u *UserProfile) Add() error {
	return db.Create(u).Error
}

func (u *UserProfile) Read() error {
	if u.ID == 0 {
		return errors.New("必须指定user_profile的id")
	}
	return db.Where("id = ?", u.ID).Find(u).Error
}

//Update .
func (u *UserProfile) Update(fields map[string]interface{}) error {
	fields["update_at"] = time.Now().Unix()
	return db.Model(u).Updates(fields).Error
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

//AddBalance .
func (u *UserProfile) AddBalance(amount int, trans *gorm.DB) error {
	if trans != nil {
		return trans.Model(u).Update("balance", gorm.Expr("balance + ?", amount)).Error
	}
	return db.Model(u).Update("balance", gorm.Expr("balance + ?", amount)).Error
}

//AllocateFund 划款
func (u *UserProfile) AllocateFund(to, invitation *UserProfile, amount, inviteAmount uint64, trans *gorm.DB) error {
	if u.Balance < amount { //检查余额
		return errors.New("用户余额不足，扣款(" + strconv.FormatUint(amount, 64) + ")失败,余额:" + strconv.FormatUint(u.Balance, 64))
	}

	if err := u.AddBalance(-int(amount), trans); err != nil { //扣款
		return errors.New("发起人扣款失败\t" + err.Error())
	}

	if err := to.AddBalance(int(amount), trans); err != nil { //增加余额
		return errors.New("接受人增加余额失败\t" + err.Error())
	}

	if invitation != nil && uint64(amount*3/10) > inviteAmount {
		if err := invitation.AddBalance(int(inviteAmount), trans); err != nil { //增加余额
			return errors.New("邀请人增加余额失败\t" + err.Error())
		}
	}

	return nil
}
