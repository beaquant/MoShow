package models

import "github.com/jinzhu/gorm"

const (
	//WithdrawStatusApply 提现申请中
	WithdrawStatusApply = iota
	//WithdrawStatusDone 提现成功
	WithdrawStatusDone
)

//Withdraw .
type Withdraw struct {
	ID       uint64 `json:"id" gorm:"column:id;primary_key"`
	UserID   uint64 `json:"user_id" gorm:"column:user_id"`
	Amount   uint64 `json:"amount" gorm:"column:amount"`
	Status   int    `json:"status" gorm:"column:status"`
	CreateAt int64  `json:"create_at" gorm:"column:create_at"`
}

//TableName .
func (Withdraw) TableName() string {
	return "withdraw"
}

//Add .
func (w *Withdraw) Add(trans *gorm.DB) error {
	if trans != nil {
		db = trans
	}

	return db.Create(w).Error
}