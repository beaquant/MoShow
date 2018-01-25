package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Guest .
type Guest struct {
	ID      uint64    `json:"id" gorm:"column:id;primary_key"`
	UserID  uint64    `json:"user_id" gorm:"column:user_id"`
	GuestID uint64    `json:"guest_id" gorm:"column:guest_id"`
	Time    time.Time `json:"time" gorm:"column:time"`
	Count   uint64    `json:"count" gorm:"column:count"`
}

//TableName .
func (Guest) TableName() string {
	return "guest"
}

//AddView .
func (g *Guest) AddView(uid, guest uint64) error {
	var count int
	if err := db.Model(g).Where("user_id = ? and guest_id = ?", uid, guest).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return db.Model(g).Where("user_id = ? and guest_id = ?", uid, guest).Updates(map[string]interface{}{"count": gorm.Expr("count + ?", 1), "time": time.Now()}).Error
	}

	gst := &Guest{UserID: uid, GuestID: guest, Time: time.Now(), Count: 1}
	return db.Model(g).Create(gst).Error
}
