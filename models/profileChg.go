package models

import (
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
	ID                     uint64 `json:"user_id" gorm:"column:id;primary_key" description:"用户ID"`
	Alias                  string `json:"alias" gorm:"column:alias"`
	AliasCheckStatus       int    `json:"alias_check" gorm:"column:alias_check"`
	Description            string `json:"description" gorm:"column:description"`
	DescriptionCheckStatus int    `json:"description_check" gorm:"column:description_check"`
	CoverPic               string `json:"cover_pic" gorm:"column:cover_pic"`
	CoverPicCheckStatus    int    `json:"cover_pic_check" gorm:"column:cover_pic_check"`
	Video                  string `json:"video" gorm:"column:video"`
	VideoCheckStatus       int    `json:"video_check" gorm:"column:video_check"`
	UpdateAt               int64  `json:"update_at" gorm:"column:update_at"`
}

//TableName .
func (ProfileChg) TableName() string {
	return "profile_chg"
}

//ReadOrCreate .
func (p *ProfileChg) ReadOrCreate(trans *gorm.DB) (err error) {
	var pl []ProfileChg

	if trans == nil {
		trans = db
	}

	err = trans.Where("id = ?", p.ID).Find(&pl).Error
	if err != nil {
		return
	}

	if pl != nil && len(pl) > 0 {
		*p = pl[0]
	} else {
		err = trans.Create(p).Error
	}

	return
}

//Update .
func (p *ProfileChg) Update(fields map[string]interface{}, trans *gorm.DB) error {
	fields["update_at"] = time.Now().Unix()
	if trans != nil {
		return trans.Model(p).Updates(fields).Error
	}
	return db.Model(p).Updates(fields).Error
}
