package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	//CheckStatusUncheck 未审核
	CheckStatusUncheck = iota
	//CheckStatusReject 驳回
	CheckStatusReject
	//CheckStatusPass 通过
	CheckStatusPass
)

//ProfileChg .
type ProfileChg struct {
	ID                  uint64 `json:"user_id" gorm:"column:id;primary_key" description:"用户ID"`
	CoverPic            string `json:"cover_pic" gorm:"column:cover_pic"`
	CoverPicCheckStatus int    `json:"cover_pic_check" gorm:"column:cover_pic_check"`
	Video               string `json:"video" gorm:"column:video"`
	VideoCheckStatus    int    `json:"video_check" gorm:"column:video_check"`
	UpdateAt            int64  `json:"update_at" gorm:"column:update_at"`
}

//TableName .
func (ProfileChg) TableName() string {
	return "profile_chg"
}

//ReadOrCreate .
func (p *ProfileChg) ReadOrCreate(trans *gorm.DB) (err error) {
	if p.ID == 0 {
		return errors.New("必须指定用户ID")
	}

	var pf ProfileChg
	if trans == nil {
		trans = db
	}
	if err = trans.Where("id = ?", p.ID).Find(&pf).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}
		err = trans.Create(p).Error
	}

	return
}

//Update .
func (p *ProfileChg) Update(fields map[string]interface{}, trans *gorm.DB) error {
	fields["update_at"] = time.Now().Unix()
	if len(fields) == 1 {
		return nil
	}

	if trans != nil {
		return trans.Model(p).Updates(fields).Error
	}
	return db.Model(p).Updates(fields).Error
}
