package models

import (
	"errors"
	"time"

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
	//BalanceChgTypeWithDraw 收益提现
	BalanceChgTypeWithDraw
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
	GiftInfo Gift   `json:"gift_info"`
}

//VideoChgInfo .
type VideoChgInfo struct {
	DialID   uint64 `json:"dial_id"`
	TimeLong uint64 `json:"time_long"`
	Price    uint64 `json:"price"`
}

//MessageOrVideoChgInfo .
type MessageOrVideoChgInfo struct {
	URL      string `json:"url"`
	TargetID uint64 `json:"target_id"`
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
	if b.Time == 0 {
		b.Time = time.Now().Unix()
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

//IsVideoPayed .
func (b *BalanceChg) IsVideoPayed(uri string, tid uint64) (bool, error) {
	var count int
	err := db.Where("user_id = ?", b.UserID).Where("chg_type = ?", BalanceChgTypeVideoView).Where(`chg_info ->>'$.target_id' = ?`, tid).Where(`chg_info ->>'$.url' = ?`, uri).Count(&count).Error
	if count > 0 {
		return true, err
	}
	return false, err
}

//GetIncomeChgs .
func (b *BalanceChg) GetIncomeChgs(limit, skip int) ([]BalanceChg, error) {
	if limit == 0 {
		limit = 20
	}

	var lst []BalanceChg
	return lst, db.Where("user_id = ?", b.UserID).Where("chg_type in (?)", []int{BalanceChgTypeRecharge, BalanceChgTypeInvitationRechargeIncome, BalanceChgTypeInvitationIncome}).Order("time desc").Limit(limit).Offset(skip).Find(&lst).Error
}

//GetPaymentChgs .
func (b *BalanceChg) GetPaymentChgs(limit, skip int) ([]BalanceChg, error) {
	if limit == 0 {
		limit = 20
	}

	var lst []BalanceChg
	return lst, db.Where("user_id = ?", b.UserID).Where("chg_type in (?)", []int{BalanceChgTypeGift, BalanceChgTypeVideo, BalanceChgTypeMessage, BalanceChgTypeVideoView}).Order("time desc").Limit(limit).Offset(skip).Find(&lst).Error
}
