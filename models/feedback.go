package models

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
}

//TableName .
func (FeedBack) TableName() string {
	return "feedback"
}

//Add .
func (f *FeedBack) Add() error {
	return db.Model(f).Create(f).Error
}
