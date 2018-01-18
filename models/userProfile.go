package models

import (
	"time"
)

const (
	//GenderMan .
	GenderMan = iota
	//GenderWoman .
	GenderWoman
)

//UserProfile .
type UserProfile struct {
	ID          uint64    `json:"user_id" gorm:"column:id;primary_key"`
	Alias       string    `json:"alias" gorm:"column:alias"`
	Gender      int       `json:"gender" gorm:"column:gender"`
	Description string    `json:"description" gorm:"column:description"`
	Birthday    time.Time `json:"birthday" gorm:"column:birthday"`
	Location    string    `json:"location" gorm:"column:location"`
}

//TableName .
func (UserProfile) TableName() string {
	return "users"
}

//Add .
func (u *UserProfile) Add() error {
	return db.Create(u).Error
}

func (u *UserProfile) Read() error {
	return db.Where("id = ?", u.ID).Find(u).Error
}

//Update .
func (u *UserProfile) Update() {

}
