package models

import "time"

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
	//AcctStatusDisable 注销账号
	AcctStatusDisable
	//AcctStatusShield 屏蔽账号
	AcctStatusShield
)

//User .
type User struct {
	ID         uint64    `json:"user_id" gorm:"column:id;primary_key"`
	UserName   string    `json:"user_name" gorm:"column:name"`
	AcctType   int       `json:"acct_type" gorm:"column:acct_type"`
	AcctStatus int       `json:"acct_status" gorm:"column:acct_status"`
	CreatedAt  time.Time `json:"create_at" gorm:"column:create_at"`
}

//TableName .
func (User) TableName() string {
	return "users"
}

//Add .
func (u *User) Add() error {
	return db.Create(u).Error
}

//ReadFromUserName .
func (u *User) ReadFromUserName() error {
	return db.Where("name = ?", u.UserName).Find(u).Error
}
