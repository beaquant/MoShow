package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
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
	Balance     int       `json:"balance" gorm:"column:balance"`
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
	if u.ID == 0 {
		return errors.New("必须只能user_profile的id")
	}
	return db.Where("id = ?", u.ID).Find(u).Error
}

//Update .
func (u *UserProfile) Update(col ...string) {

}

//AddBalance .
func (u *UserProfile) AddBalance(amount int, trans *gorm.DB) {
	if trans != nil {
		trans.Model(u).Update("balance", gorm.Expr("balance +", amount))
	} else {
		db.Model(u).Update("balance", gorm.Expr("balance +", amount))
	}
}
