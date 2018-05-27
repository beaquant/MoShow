package models

//Active .
type Active struct {
	ID        uint64 `json:"id" gorm:"column:id;primary_key"`
	UserID    uint64 `json:"user_id" gorm:"column:user_id"`
	Type      int    `json:"active_type" gorm:"column:active_type"`
	Content   string `json:"content" gorm:"column:content"`
	DelayTime int    `json:"delay_time" gorm:"column:delay_time"`
}

//TableName .
func (Active) TableName() string {
	return "active"
}
