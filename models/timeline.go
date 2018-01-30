package models

import "time"

//TimelineUser .
type TimelineUser struct {
	ID          uint64    `json:"user_id" gorm:"column:id;primary_key"`
	Alias       string    `json:"alias" gorm:"column:alias"`
	Gender      int       `json:"gender" gorm:"column:gender"`
	CoverPic    string    `json:"-" gorm:"column:cover_pic"`
	Description string    `json:"description" gorm:"column:description"`
	Birthday    time.Time `json:"birthday" gorm:"column:birthday"`
	Location    string    `json:"location" gorm:"column:location"`
	Balance     uint64    `json:"balance" gorm:"column:balance"`
	Price       uint64    `json:"price" gorm:"column:price"`
	Duration    uint64    `json:"duration" gorm:"column:duration"`
}

//TableName .
func (TimelineUser) TableName() string {
	return "time_line"
}
