package models

//UserExtra 用户附加信息
type UserExtra struct {
	ID          uint64 `json:"user_id" gorm:"column:id;primary_key" description:"用户ID"`
	GiftHistory string `json:"-" gorm:"column:gift_his" description:"收到的礼物"`
}

//GiftHisInfo .
type GiftHisInfo struct {
	GiftInfo Gift
	Count    uint64
}

//TableName .
func (UserExtra) TableName() string {
	return "user_extra"
}

//GetGiftHis .
func (u *UserExtra) GetGiftHis() {
	// gftHist := make(map[uint64]GiftHisInfo)
}
