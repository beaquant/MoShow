package models

//Dial .
type Dial struct {
	ID         uint64 `json:"id" gorm:"column:id;primary_key"`
	FromUserID uint64 `json:"from_user_id" gorm:"column:from_user_id"`
	ToUserID   uint64 `json:"to_user_id" gorm:"column:to_user_id"`
	Duration   int    `json:"duration" gorm:"column:duration"`
	CreateAt   int64  `json:"create_at" gorm:"column:create_at"`
	Success    bool   `json:"success" gorm:"column:success"`
}

//TableName .
func (Dial) TableName() string {
	return "dial"
}

//Add .
func (d *Dial) Add() error {
	return db.Model(d).Create(d).Error
}
