package models

import "time"

//Guest .
type Guest struct {
	ID      uint64    `json:"id" gorm:"column:id;primary_key"`
	UserID  uint64    `json:"user_id" gorm:"column:user_id"`
	GuestID uint64    `json:"guest_id" gorm:"column:guest_id"`
	Time    time.Time `json:"time" gorm:"column:time"`
}

//TableName .
func (Guest) TableName() string {
	return "guest"
}
