package models

import (
	"github.com/jinzhu/gorm"
)

const (
	//PayTypeAlipay 支付宝
	PayTypeAlipay = iota
	//PayTypeWechatpay 微信支付
	PayTypeWechatpay
	//PayTypeApplepurchase 苹果内购
	PayTypeApplepurchase
)

//Order .
type Order struct {
	ID          uint64  `json:"id" gorm:"column:id;primary_key"`
	UserID      uint64  `json:"user_id" gorm:"column:user_id"`
	Amount      float64 `json:"amount" gorm:"column:amount"`
	CoinCount   uint64  `json:"coin_count" gorm:"column:coin_count"`
	Success     bool    `json:"success" gorm:"column:success"`
	PayType     int     `json:"pay_type" gorm:"column:pay_type"`
	CreateAt    int64   `json:"create_at" gorm:"column:create_at"`
	PayTime     int64   `json:"pay_time" gorm:"column:pay_time"`
	PayInfo     string  `json:"pay_info" gorm:"column:pay_info"`
	ProductInfo string  `json:"product_info" gorm:"column:product_info"`
}

//TableName .
func (Order) TableName() string {
	return "order"
}

//Add .
func (o *Order) Add(trans *gorm.DB) error {
	if trans != nil {
		return trans.Model(o).Create(o).Error
	}
	return db.Model(o).Create(o).Error
}

//Read .
func (o *Order) Read() error {
	return db.Where("id = ?", o.ID).Find(o).Error
}

//Update .
func (o *Order) Update(fields map[string]interface{}, trans *gorm.DB) error {
	if trans != nil {
		return trans.Model(o).Updates(fields).Error
	}
	return db.Model(o).Updates(fields).Error
}
