package models

import (
	"MoShow/utils"
	"errors"
	"time"
)

const (
	//FeedBackTypeSuggestion 建议
	FeedBackTypeSuggestion = iota
	//FeedBackTypeReport 举报
	FeedBackTypeReport
)

//FeedBack .
type FeedBack struct {
	ID      uint64 `json:"id" gorm:"column:id;primary_key"`
	UserID  uint64 `json:"user_id" gorm:"column:user_id"`
	Type    int    `json:"type" gorm:"column:type"`
	Content string `json:"content" gorm:"column:content"`
	Time    int64  `json:"time" gorm:"column:time"`
}

//FeedBackReport 用户举报
type FeedBackReport struct {
	FeedBackSuggestion
	TgUserID uint64 `json:"target_user_id"`
	Cate     string `json:"cate"`
}

//FeedBackSuggestion 意见反馈
type FeedBackSuggestion struct {
	Content string `json:"feedback_content"`
	Img     string `json:"img"`
}

//TableName .
func (FeedBack) TableName() string {
	return "feedback"
}

//AddReport .
func (f *FeedBack) AddReport(r *FeedBackReport) (err error) {
	if r.TgUserID == 0 {
		return errors.New("被举报人的ID不能为0")
	}

	f.Type = FeedBackTypeReport
	f.Time = time.Now().Unix()
	if f.Content, err = utils.JSONMarshalToString(r); err != nil {
		return err
	}

	return db.Model(f).Create(f).Error
}

//AddSuggestion .
func (f *FeedBack) AddSuggestion(s *FeedBackSuggestion) (err error) {
	f.Type = FeedBackTypeSuggestion
	f.Time = time.Now().Unix()
	if f.Content, err = utils.JSONMarshalToString(s); err != nil {
		return err
	}

	return db.Model(f).Create(f).Error
}
