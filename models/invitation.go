package models

import "time"

//Invitation .
type Invitation struct {
	ID         uint64    `json:"id" gorm:"column:id;primary_key"`
	FromUserID uint64    `json:"from_user_id" gorm:"column:from_user_id"`
	ToUserID   uint64    `json:"to_user_id" gorm:"column:to_user_id"`
	Time       time.Time `json:"time" gorm:"column:time"`
}

//TableName .
func (Invitation) TableName() string {
	return "invitation"
}
