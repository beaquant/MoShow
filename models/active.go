package models

const (
	//ActiveTypeMessage 文字促活消息
	ActiveTypeMessage = iota
	//ActiveTypeImage 图片促活消息
	ActiveTypeImage
	//ActiveTypeVoice 语音促活消息
	ActiveTypeVoice
	//ActiveTypeVideo 视频促活消息
	ActiveTypeVideo
)

//Active .
type Active struct {
	ID        uint64 `json:"id" gorm:"column:id;primary_key"`
	UserID    uint64 `json:"user_id" gorm:"column:user_id"`
	Type      int    `json:"active_type" gorm:"column:active_type"`
	Content   string `json:"content" gorm:"column:content"`
	DelayTime int64  `json:"delay_time" gorm:"column:delay_time"`
}

//ActiveDetail .
type ActiveDetail struct {
	Message  string `json:"msg"`
	FileURL  string `json:"url"`
	Duration uint   `json:""`
}

//TableName .
func (Active) TableName() string {
	return "active"
}

//GetActive .
func (Active) GetActive() (act []Active, err error) {
	err = db.Find(&act).Error
	return
}
