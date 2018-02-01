package models

//账号类型
const (
	//AcctTypeTelephone 手机号账户
	AcctTypeTelephone = iota
	//AcctTypeWechat 微信登陆账户
	AcctTypeWechat
)

//账号状态
const (
	//AcctStatusNormal 正常账号
	AcctStatusNormal = iota
	//AcctStatusDeleted 注销账号
	AcctStatusDeleted
	//AcctStatusShield 屏蔽账号
	AcctStatusShield
)

//User .
type User struct {
	ID          uint64 `json:"user_id" gorm:"column:id;primary_key" description:"用户ID"`
	PhoneNumber string `json:"phone_number" gorm:"column:phone_number" description:"手机号"`
	WeChatID    string `json:"wechat_id" gorm:"column:wechat_id" description:"微信ID"`
	AcctType    int    `json:"acct_type" gorm:"column:acct_type" description:"账号类型"`
	AcctStatus  int    `json:"acct_status" gorm:"column:acct_status" description:"账号状态"`
	CreatedAt   int64  `json:"create_at" gorm:"column:create_at" description:"注册时间"`
	InvitedBy   uint64 `json:"invited_by" gorm:"column:invited_by" description:"邀请人ID"`
}

//TableName .
func (User) TableName() string {
	return "users"
}

//Add .
func (u *User) Add() error {
	return db.Create(u).Error
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
