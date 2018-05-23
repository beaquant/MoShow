package models

import (
	"MoShow/utils"

	"github.com/jinzhu/gorm"
)

//账号类型
const (
	//AcctTypeTelephone 手机号账户
	AcctTypeTelephone = iota
	//AcctTypeWechat 微信登陆账户
	AcctTypeWechat
	//AcctTypeInternal 内部账户
	AcctTypeInternal
)

//账号状态
const (
	//AcctStatusNormal 正常账号
	AcctStatusNormal = iota
	//AcctStatusShield 封禁账号
	AcctStatusShield
)

//User .
type User struct {
	ID            uint64 `json:"user_id" gorm:"column:id;primary_key" description:"用户ID"`
	PhoneNumber   string `json:"phone_number" gorm:"column:phone_number" description:"手机号"`
	WeChatID      string `json:"wechat_id" gorm:"column:wechat_id" description:"微信ID"`
	AcctType      int    `json:"acct_type" gorm:"column:acct_type" description:"账号类型"`
	AcctStatus    int    `json:"acct_status" gorm:"column:acct_status" description:"账号状态"`
	CreatedAt     int64  `json:"create_at" gorm:"column:create_at" description:"注册时间"`
	InvitedBy     uint64 `json:"invited_by" gorm:"column:invited_by" description:"邀请人ID"`
	InvitedAward  uint64 `json:"-" gorm:"column:invited_award" description:"给邀请人的奖励"`
	LastLoginInfo string `json:"last_login_info" gorm:"column:last_login_info" description:"最近一次登录信息"`
}

//UserLoginInfo .
type UserLoginInfo struct {
	UserAgent string `json:"user_agent"`
	IPAddress string `json:"ip_address"`
	Time      int64  `json:"time"`
}

//TableName .
func (User) TableName() string {
	return "users"
}

//Add .
func (u *User) Add(trans *gorm.DB) error {
	if trans != nil {
		return trans.Create(u).Error
	}
	return db.Create(u).Error
}

func (u *User) Read() error {
	return db.Find(u).Error
}

//UpdateLoginInfo .
func (u *User) UpdateLoginInfo(uli *UserLoginInfo) error {
	str, err := utils.JSONMarshalToString(uli)
	if err != nil {
		return err
	}
	return db.Model(u).Update("last_login_info", str).Error
}

//ReadFromPhoneNumber .
func (u *User) ReadFromPhoneNumber() (err error) {
	var ul []User
	err = db.Where("phone_number = ?", u.PhoneNumber).Find(&ul).Error
	if ul != nil && len(ul) > 0 {
		*u = ul[0]
	}

	return
}

//ReadFromWechatID .
func (u *User) ReadFromWechatID() (err error) {
	var ul []User
	err = db.Where("wechat_id = ?", u.WeChatID).Find(&ul).Error
	if ul != nil && len(ul) > 0 {
		*u = ul[0]
	}

	return
}

//GetInternalAcct .
func (u *User) GetInternalAcct() ([]User, error) {
	var uArr []User
	return uArr, db.Where("acct_type = ?", AcctTypeInternal).Find(&uArr).Error
}

//AddAward .
func (u *User) AddAward(count uint64, trans *gorm.DB) error {
	if trans == nil {
		trans = db
	}
	return trans.Model(u).Update("invited_award", gorm.Expr("invited_award + ?", count)).Error
}

//GetRegistTimeAndAward .
func (u *User) GetRegistTimeAndAward() error {
	return db.Select("create_at,invited_award").Find(u).Error
}
