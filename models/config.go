package models

import (
	"MoShow/utils"
	"time"
)

const (
	configTypeGift       = "gifts"
	configTypeProduct    = "products"
	configTypeIncomeRate = "income_rate"
	configTypeBanner     = "banner"
)

var (
	giftList    map[string]Gift
	productList map[string]Product
	updateTime  = make(map[string]time.Time)
	incomeRate  *IncomeRate
	banners     []Banner
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
	Price    uint64 `json:"price"`
	ImgURL   string `json:"img"`
}

//Product .
type Product struct {
	ProductName string  `json:"name"`
	Code        string  `json:"code"`
	Price       float64 `json:"price"`
	CoinCount   uint64  `json:"coin_count"`
}

//IncomeRate 分成比例
type IncomeRate struct {
	InviteRechargeRate float64 `json:"invite_recharge_rate"` //被邀请人充值分成率
	InviteIncomegeRate float64 `json:"invite_income_rate"`   //被邀请人收益分成
	IncomeFee          float64 `json:"income_fee"`           //收益手续费
}

const (
	//BannerTypeImg 纯图片banner
	BannerTypeImg = iota
	//BannerTypeLink 链接跳转banner
	BannerTypeLink
)

//Banner 首页banner
type Banner struct {
	Image      string
	URL        string
	BannerType int
}

//TableName .
func (Config) TableName() string {
	return "config"
}

//GetCommonGiftInfo 获取礼物列表
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

//GetProductInfo 获取商品列表
func (c *Config) GetProductInfo() (map[string]Product, error) {
	if tm, ok := updateTime[configTypeProduct]; productList == nil || !ok || tm.Add(time.Minute*5).Before(time.Now()) {
		if err := db.Debug().Where("conf_key = ?", configTypeProduct).First(c).Error; err != nil {
			return nil, err
		}

		pf := make(map[string]Product)
		if err := utils.JSONUnMarshal(c.Value, &pf); err != nil {
			return nil, err
		}

		productList = pf
		updateTime[configTypeProduct] = time.Now()
		return pf, nil
	}
	return productList, nil
}

//GetIncomeRate 获取分成比例
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

//GetBanners 获取banner
func (c *Config) GetBanners() {
	if tm, ok := updateTime[configTypeBanner]; banners == nil || !ok || tm.Add(time.Minute*5).Before(time.Now()) {
	}
}
