package models

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

const (
	//DialStatusFail 通话状态,失败
	DialStatusFail = iota
	//DialStatusSuccess 成功
	DialStatusSuccess
	//DialStatusException 异常
	DialStatusException
)

//Dial .
type Dial struct {
	ID         uint64 `json:"id" gorm:"column:id;primary_key"`
	FromUserID uint64 `json:"from_user_id" gorm:"column:from_user_id"`
	ToUserID   uint64 `json:"to_user_id" gorm:"column:to_user_id"`
	Duration   uint64 `json:"duration" gorm:"column:duration"`
	CreateAt   int64  `json:"create_at" gorm:"column:create_at"`
	Status     int    `json:"status" gorm:"column:status"`
	Tag        string `json:"tag" gorm:"column:tag"`
	Clearing   string `json:"-" gorm:"column:clearing"`
}

//DialTag .
type DialTag struct {
	ErrorMsg []string
}

//ClearingInfo .
type ClearingInfo struct {
	NIMChannelID uint64 //网易云信房间ID
	Cost         uint64 `json:"cost" description:"用户花费"`
	Income       uint64 `json:"income,omitempty" description:"主播收益"`
	Timelong     uint64 `json:"timelong" description:"聊天时长"`
}

//TableName .
func (Dial) TableName() string {
	return "dial"
}

//Add .
func (d *Dial) Add() error {
	if len(d.Tag) == 0 {
		d.Tag = "null"
	}

	if len(d.Clearing) == 0 {
		d.Clearing = "null"
	}

	return db.Model(d).Create(d).Error
}

//Update .
func (d *Dial) Update(fields map[string]interface{}, trans *gorm.DB) error {
	if trans == nil {
		trans = db
	}

	return trans.Model(d).Updates(fields).Error
}

//Read .
func (d *Dial) Read() error {
	return db.Where("id = ?", d.ID).Find(d).Error
}

//GetDialList .
func (d *Dial) GetDialList(uid uint64, limit, skip int) ([]Dial, error) {
	if limit == 0 {
		limit = 20
	}

	var lst []Dial
	return lst, db.Where("from_user_id = ?", uid).Find(&lst).Order("create_at desc").Limit(limit).Offset(skip).Error
}

//Del .
func (d *Dial) Del() error {
	return db.Delete(d).Error
}

//GetToalDuration .
func (d *Dial) GetToalDuration(uid uint64) (int64, error) {
	row := db.Table(d.TableName()).Where("status = ? ", DialStatusSuccess).Where("from_user_id = ? or to_user_id = ?", uid, uid).Select("sum(duration) as duration").Row()
	var duration sql.NullInt64
	if err := row.Scan(&duration); err != nil {
		return 0, err
	}

	if duration.Valid {
		return duration.Int64, nil
	}
	return 0, nil
}
