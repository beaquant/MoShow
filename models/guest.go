package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Guest .
type Guest struct {
	ID      uint64 `json:"id" gorm:"column:id;primary_key"`
	UserID  uint64 `json:"user_id" gorm:"column:user_id"`
	GuestID uint64 `json:"guest_id" gorm:"column:guest_id"`
	Time    int64  `json:"time" gorm:"column:time"`
	Count   uint64 `json:"count" gorm:"column:count"`
}

//TableName .
func (Guest) TableName() string {
	return "guest"
}

//AddView .
func (g *Guest) AddView(uid, guest uint64) error {
	var gg []Guest
	if err := db.Where("user_id = ? and guest_id = ?", uid, guest).Find(&gg).Error; err != nil {
		return err
	}

	if gg != nil && len(gg) > 0 {
		return db.Model(gg[0]).Updates(map[string]interface{}{"count": gorm.Expr("count + ?", 1), "time": time.Now().Unix()}).Error
	}

	gst := &Guest{UserID: uid, GuestID: guest, Time: time.Now().Unix(), Count: 1}
	return db.Model(g).Create(gst).Error
}

//GetGuestList 获取指定用户的访客列表
func (g *Guest) GetGuestList(uid uint64, limit, skip int) ([]Guest, error) {
	if limit == 0 {
		limit = 20
	}

	var gg []Guest
	return gg, db.Where("user_id = ?", uid).Limit(limit).Offset(skip).Order("time desc").Find(&gg).Error
}
