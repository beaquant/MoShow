package models

import "time"

const (
	//BalanceChgTypeRecharge 充值增加余额
	BalanceChgTypeRecharge = iota
	//BalanceChgTypeInvitationIncome 邀请分成
	BalanceChgTypeInvitationIncome
	//BalanceChgTypeSendGift 赠送礼物消费
	BalanceChgTypeSendGift
	//BalanceChgTypeReceiveGift 收到礼物
	BalanceChgTypeReceiveGift
)

//BalanceChg .
type BalanceChg struct {
	ID      uint64    `json:"id" gorm:"column:id;primary_key"`
	UserID  uint64    `json:"user_id" gorm:"column:user_id"`
	ChgType int       `json:"chg_type" gorm:"column:chg_type"`
	ChgInfo string    `json:"chg_info" gorm:"column:chg_info"`
	Amount  int       `json:"amount" gorm:"column:amount"`
	Time    time.Time `json:"time" gorm:"column:time"`
}

//TableName .
func (BalanceChg) TableName() string {
	return "balance_chg"
}
