package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

const (
	//BalanceChgTypeRecharge 充值增加余额
	BalanceChgTypeRecharge = iota
	//BalanceChgTypeSendGift 赠送礼物消费
	BalanceChgTypeSendGift
	//BalanceChgTypeReceiveGift 收到礼物
	BalanceChgTypeReceiveGift
	//BalanceChgTypeInvitationRechargeIncome 邀请用户充值分成
	BalanceChgTypeInvitationRechargeIncome
	//BalanceChgTypeInvitationIncome 邀请用户收益分成
	BalanceChgTypeInvitationIncome
)

//BalanceChg .
type BalanceChg struct {
	ID         uint64 `json:"id" gorm:"column:id;primary_key"`
	UserID     uint64 `json:"user_id" gorm:"column:user_id"`
	FromUserID uint64 `json:"from_user_id" gorm:"column:from_user_id"`
	ChgType    int    `json:"chg_type" gorm:"column:chg_type"`
	ChgInfo    string `json:"chg_info" gorm:"column:chg_info"`
	Amount     int    `json:"amount" gorm:"column:amount"`
	Time       int64  `json:"time" gorm:"column:time"`
}

//GiftChgInfo .
type GiftChgInfo struct {
	Count    uint64
	GiftInfo *Gift
}

//VideoChgInfo .
type VideoChgInfo struct {
	TimeLong uint64
}

//BalanceChgInfo .
type BalanceChgInfo struct {
	Content string `json:"content"`
}

//TableName .
func (BalanceChg) TableName() string {
	return "balance_chg"
}

//Add .
func (b *BalanceChg) Add(trans *gorm.DB) error {
	if b.UserID == 0 {
		return errors.New("余额变动记录必须指定用户ID")
	}

	if trans != nil {
		return trans.Create(b).Error
	}
	return db.Create(b).Error
}

//AddChg .
func (b *BalanceChg) AddChg(trans *gorm.DB, chg ...*BalanceChg) error {
	if trans == nil {
		return errors.New("事务不能为空")
	}

	for index := range chg {
		if chg[index] == nil {
			break
		}

		if err := chg[index].Add(trans); err != nil {
			return err
		}
	}
	return nil
}
