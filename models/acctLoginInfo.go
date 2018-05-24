package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

//AcctLoginInfo 账号登陆信息
type AcctLoginInfo struct {
	ID        uint64 `json:"id" gorm:"column:id;primary_key"`
	UserID    uint64 `json:"user_id" gorm:"column:user_id"`
	IPAddress string `json:"ip_address" gorm:"column:ip_address"`
	Agent     string `json:"agent" gorm:"column:agent"`
	Time      int64  `json:"time" gorm:"column:time"`
	IPInfo    string `json:"ip_info" gorm:"column:ip_info"`
}

//TableName .
func (AcctLoginInfo) TableName() string {
	return "acct_login_info"
}

//Add .
func (a *AcctLoginInfo) Add(trans *gorm.DB) error {
	if a.UserID == 0 {
		return errors.New("账号登陆信息必须指定用户ID")
	}
	if a.Time == 0 {
		a.Time = time.Now().Unix()
	}

	if trans == nil {
		trans = db
	}

	return trans.Create(a).Error
}
