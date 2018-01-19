package models

import (
	"time"
)

//Order .
type Order struct {
	ID       uint64    `json:"id" gorm:"column:id;primary_key"`
	UserID   uint64    `json:"user_id" gorm:"column:user_id"`
	Amount   float32   `json:"amount" gorm:"column:amount"`
	Success  bool      `json:"success" gorm:"column:success"`
	PayType  int       `json:"pay_type" gorm:"column:pay_type"`
	CreateAt time.Time `json:"create_at" gorm:"column:create_at"`
	PayTime  time.Time `json:"pay_time" gorm:"column:pay_time"`
	PayInfo  string    `json:"pay_info" gorm:"column:pay_info"`
}

//TableName .
func (Order) TableName() string {
	return "order"
}
