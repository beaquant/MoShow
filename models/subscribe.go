package models

import (
	"MoShow/utils"
	"errors"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
)

//Subscribe 关注与被关注信息
type Subscribe struct {
	ID        uint64 `json:"user_id" gorm:"column:id;primary_key" description:"用户ID"`
	Followers string `json:"-" gorm:"column:follower" description:"关注者"`
	Following string `json:"-" gorm:"column:following" description:"正在关注"`
}

//FollowInfo .
type FollowInfo struct {
	FollowTime int64
}

//TableName .
func (Subscribe) TableName() string {
	return "subscribe"
}

//Add .
func (s *Subscribe) Add(trans *gorm.DB) error {
	if trans == nil {
		trans = db
	}

	if len(s.Followers) == 0 {
		s.Followers = "{}"
	}

	if len(s.Following) == 0 {
		s.Following = "{}"
	}

	return trans.Create(s).Error
}

//ReadOrCreate .
func (s *Subscribe) Read(trans *gorm.DB) error {
	if s.ID == 0 {
		return errors.New("必须指定用户ID")
	}

	if trans == nil {
		trans = db
	}
	return trans.Find(s).Error
}

//AddFollow 添加关注
func (s *Subscribe) AddFollow(id uint64) error {
	idStr := strconv.FormatUint(id, 10)
	fis, _ := utils.JSONMarshalToString(&FollowInfo{FollowTime: time.Now().Unix()})

	trans := db.Begin()
	if err := trans.Model(s).Update("following", gorm.Expr(`JSON_SET(COALESCE(following,'{}'),'$."`+idStr+`"',CAST('`+fis+`' AS JSON))`)).Error; err != nil {
		trans.Rollback()
		return err
	}

	if err := trans.Model(&Subscribe{ID: id}).Update("follower", gorm.Expr(`JSON_SET(COALESCE(follower,'{}'),'$."`+idStr+`"',CAST('`+fis+`' AS JSON))`)).Error; err != nil {
		trans.Rollback()
		return err
	}

	trans.Commit()
	return nil
}

//UnFollow 取消关注
func (s *Subscribe) UnFollow(id uint64) error {
	idStr := strconv.FormatUint(id, 10)
	trans := db.Begin()
	if err := trans.Model(s).Update("following", gorm.Expr(`JSON_REMOVE(follower,'$."`+idStr+`"')`)).Error; err != nil {
		trans.Rollback()
		return err
	}

	if err := trans.Model(&Subscribe{ID: id}).Update("follower", gorm.Expr(`JSON_REMOVE(follower,'$."`+idStr+`"')`)).Error; err != nil {
		trans.Rollback()
		return err
	}

	trans.Commit()
	return nil
}

//GetFollowers .
func (s *Subscribe) GetFollowers() map[uint64]FollowInfo {
	if len(s.Followers) > 0 {
		fl := make(map[uint64]FollowInfo)
		if err := utils.JSONUnMarshal(s.Followers, &fl); err != nil {
			beego.Error(err)
			return nil
		}
		return fl
	}
	return nil
}

//GetFollowing .
func (s *Subscribe) GetFollowing() map[uint64]FollowInfo {
	if len(s.Following) > 0 {
		fl := make(map[uint64]FollowInfo)
		if err := utils.JSONUnMarshal(s.Following, &fl); err != nil {
			beego.Error(err)
			return nil
		}
		return fl
	}
	return nil
}
