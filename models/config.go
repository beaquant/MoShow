package models

import (
	"MoShow/utils"
)

const configTypeGift = "gifts"

//Config .
type Config struct {
	ID    uint64 `json:"id" gorm:"column:id;primary_key"`
	Key   string `json:"key" gorm:"column:key"`
	Value string `json:"val" gorm:"column:val"`
}

//TableName .
func (Config) TableName() string {
	return "config"
}

//GetCommonGiftInfo .
func (c *Config) GetCommonGiftInfo() (*CommonGiftInfo, error) {
	if err := db.Model(c).Where("key = ?", configTypeGift).First(c).Error; err != nil {
		return nil, err
	}

	gf := &CommonGiftInfo{}
	if err := utils.JSONUnMarshal(c.Value, gf); err != nil {
		return nil, err
	}
	return gf, nil
}

//CommonGiftInfo .
type CommonGiftInfo struct {
	GiftList []Gift `json:"gift_list"`
}

//Gift .
type Gift struct {
	GiftName string `json:"name"`
	Code     string `json:"code"`
	Price    uint   `json:"price"`
	ImgURL   string `json:"img"`
}
