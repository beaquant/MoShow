package models

//Order .
type Order struct {
	ID        uint64  `json:"id" gorm:"column:id;primary_key"`
	UserID    uint64  `json:"user_id" gorm:"column:user_id"`
	Amount    float32 `json:"amount" gorm:"column:amount"`
	CoinCount uint64  `json:"coin_count" gorm:"column:coin_count"`
	Success   bool    `json:"success" gorm:"column:success"`
	PayType   int     `json:"pay_type" gorm:"column:pay_type"`
	CreateAt  int64   `json:"create_at" gorm:"column:create_at"`
	PayTime   int64   `json:"pay_time" gorm:"column:pay_time"`
	PayInfo   string  `json:"pay_info" gorm:"column:pay_info"`
}

//TableName .
func (Order) TableName() string {
	return "order"
}

//Add .
func (o *Order) Add() error {
	return db.Model(o).Create(o).Error
}

//Read .
func (o *Order) Read() error {
	return db.Where("id = ?", o.ID).Find(o).Error
}
