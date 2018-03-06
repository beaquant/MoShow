package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

const (
	//BalanceChgTypeRecharge 充值增加余额
	BalanceChgTypeRecharge = iota
	//BalanceChgTypeGift 礼物消费
	BalanceChgTypeGift
	//BalanceChgTypeVideo 视频结算
	BalanceChgTypeVideo
	//BalanceChgTypeMessage 消息扣费
	BalanceChgTypeMessage
	//BalanceChgTypeVideoView 查看视频
	BalanceChgTypeVideoView
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
	Count    uint64 `json:"count"`
	GiftInfo *Gift  `json:"gift_info"`
}

//VideoChgInfo .
type VideoChgInfo struct {
	TimeLong uint64 `json:"time_long"`
	Price    uint64 `json:"price"`
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

//Read .
func (b *BalanceChg) Read() error {
	return db.First(b, b.ID).Error
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

//GetIncomeChgs .
func (b *BalanceChg) GetIncomeChgs(limit, skip int) ([]BalanceChg, error) {
	if limit == 0 {
		limit = 20
	}

	var lst []BalanceChg
	return lst, db.Where("user_id = ?", b.UserID).Where("chg_type in (?)", []int{BalanceChgTypeInvitationRechargeIncome, BalanceChgTypeInvitationIncome}).Find(&lst).Order("time").Limit(limit).Offset(skip).Error
}

//GetPaymentChgs .
func (b *BalanceChg) GetPaymentChgs(limit, skip int) ([]BalanceChg, error) {
	if limit == 0 {
		limit = 20
	}

	var lst []BalanceChg
	return lst, db.Where("user_id = ?", b.UserID).Where("chg_type in (?)", []int{BalanceChgTypeGift, BalanceChgTypeVideo, BalanceChgTypeMessage, BalanceChgTypeVideoView}).Find(&lst).Order("time").Limit(limit).Offset(skip).Error
}
