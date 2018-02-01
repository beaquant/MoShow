package models

import (
	"MoShow/utils"
	"time"
)

const configTypeGift = "gifts"
const configTypeIncomeRate = "income_rate"

var (
	giftList   map[string]Gift
	updateTime = make(map[string]time.Time)
	incomeRate *IncomeRate
)

//Config .
type Config struct {
	ID    uint64 `json:"id" gorm:"column:id;primary_key"`
	Key   string `json:"key" gorm:"column:conf_key"`
	Value string `json:"val" gorm:"column:val"`
}

//Gift .
type Gift struct {
	GiftName string `json:"name"`
	Code     string `json:"code"`
	Price    uint   `json:"price"`
	ImgURL   string `json:"img"`
}

//IncomeRate 分成比例
type IncomeRate struct {
	RechargeIncomeRate float64 //充值分成率
	VideoIncomeRate    float64 //视频聊天消费分成率
}

//TableName .
func (Config) TableName() string {
	return "config"
}

//GetCommonGiftInfo .
func (c *Config) GetCommonGiftInfo() (map[string]Gift, error) {
	if tm, ok := updateTime[configTypeGift]; giftList == nil || !ok || tm.Add(time.Minute*5).Before(time.Now()) {
		if err := db.Debug().Where("conf_key = ?", configTypeGift).First(c).Error; err != nil {
			return nil, err
		}

		gf := make(map[string]Gift)
		if err := utils.JSONUnMarshal(c.Value, &gf); err != nil {
			return nil, err
		}

		giftList = gf
		updateTime[configTypeGift] = time.Now()
		return gf, nil
	}
	return giftList, nil
}

//GetIncomeRate .
func (c *Config) GetIncomeRate() (*IncomeRate, error) {
	if tm, ok := updateTime[configTypeIncomeRate]; incomeRate == nil || !ok || tm.Add(time.Minute*5).Before(time.Now()) {
		if err := db.Model(c).Where("conf_key = ?", configTypeIncomeRate).First(c).Error; err != nil {
			return nil, err
		}

		var ir IncomeRate
		if err := utils.JSONUnMarshal(c.Value, &ir); err != nil {
			return nil, err
		}

		updateTime[configTypeIncomeRate] = time.Now()
		incomeRate = &ir
		return &ir, nil
	}
	return incomeRate, nil
}
