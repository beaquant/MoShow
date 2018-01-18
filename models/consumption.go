package models

import "time"

//Consumption .
type Consumption struct {
	ID         uint64    `json:"id" gorm:"column:id;primary_key"`
	FromUserID uint64    `json:"from_user_id" gorm:"column:from_user_id"`
	ToUserID   uint64    `json:"to_user_id" gorm:"column:to_user_id"`
	Time       time.Time `json:"time" gorm:"column:time"`
	GiftID     uint64    `json:"gift_id" gorm:"column:gift_id"`
	Count      int       `json:"count" gorm:"column:count"`
}

//TableName .
func (Consumption) TableName() string {
	return "consumption"
}
