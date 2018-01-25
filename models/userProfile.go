package models

import (
	"MoShow/utils"
	"errors"
	"time"

	"github.com/astaxie/beego"

	"github.com/jinzhu/gorm"
)

const (
	//GenderWoman .
	GenderWoman = iota
	//GenderMan .
	GenderMan
)

//UserProfile .
type UserProfile struct {
	ID          uint64    `json:"user_id" gorm:"column:id;primary_key"`
	Alias       string    `json:"alias" gorm:"column:alias"`
	Gender      int       `json:"gender" gorm:"column:gender"`
	CoverPic    string    `json:"-" gorm:"column:cover_pic"`
	Description string    `json:"description" gorm:"column:description"`
	Birthday    time.Time `json:"birthday" gorm:"column:birthday"`
	Location    string    `json:"location" gorm:"column:location"`
	Balance     uint64    `json:"balance" gorm:"column:balance"`
	Price       uint64    `json:"price" gorm:"column:price"`
}

//UserCoverInfo .
type UserCoverInfo struct {
	CoverPicture *Picture  `json:"cover_pic_info"`
	DesVideo     *Video    `json:"video_info"`
	Gallery      []Picture `json:"gallery"`
}

//Picture .
type Picture struct {
	ImageURL string `json:"image_url"`
	Disable  bool   `json:"disable"`
	Checked  bool   `json:"checked"`
}

//Video .
type Video struct {
	VideoURL string `json:"video_url"`
	Disable  bool   `json:"disable"`
	Checked  bool   `json:"checked"`
}

//TableName .
func (UserProfile) TableName() string {
	return "users"
}

//ToString .
func (u *UserCoverInfo) ToString() string {
	str, err := utils.JSONMarshalToString(u)
	if err != nil {
		beego.Error(err)
		return ""
	}
	return str
}

//Add .
func (u *UserProfile) Add() error {
	return db.Create(u).Error
}

func (u *UserProfile) Read() error {
	if u.ID == 0 {
		return errors.New("必须指定user_profile的id")
	}
	return db.Where("id = ?", u.ID).Find(u).Error
}

//Update .
func (u *UserProfile) Update(fields map[string]interface{}) error {
	return db.Model(u).Updates(fields).Error
}

//GetCover .
func (u *UserProfile) GetCover() *UserCoverInfo {
	if len(u.CoverPic) > 0 {
		ucp := &UserCoverInfo{}
		if err := utils.JSONUnMarshal(u.CoverPic, ucp); err != nil {
			beego.Error(err)
			return nil
		}
		return ucp
	}
	return nil
}

//AddBalance .
func (u *UserProfile) AddBalance(amount int, trans *gorm.DB) {
	if trans != nil {
		trans.Model(u).Update("balance", gorm.Expr("balance + ?", amount))
	} else {
		db.Model(u).Update("balance", gorm.Expr("balance + ?", amount))
	}
}
