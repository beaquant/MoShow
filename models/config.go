package models

import (
	"MoShow/utils"
	"time"

	"github.com/astaxie/beego"
)

const (
	configTypeGift       = "gifts"
	configTypeProduct    = "products"
	configTypeIncomeRate = "income_rate"
	configTypeCommon     = "common_config"
)

var (
	giftList     []Gift
	productList  []Product
	updateTime   = make(map[string]time.Time)
	incomeRate   *IncomeRate
	commonConfig *CommonConfig
)

//Config .
type Config struct {
	ID    uint64 `json:"id" gorm:"column:id;primary_key"`
	Key   string `json:"key" gorm:"column:conf_key"`
	Value string `json:"val" gorm:"column:val"`
}

//CommonConfig 通用配置
type CommonConfig struct {
	AnchorVideoRecord     bool            `json:"anchor_video_record"`
	UserVideoRecordbool   bool            `json:"user_video_record"`
	UserProtocol          string          `json:"user_protocal"`
	ForceUpdate           ForceUpdateInfo `json:"force_update"`
	Share                 ShareInfo       `json:"share"`
	CustomerServiceWechat string          `json:"customer_service_wechat"`
	CheckStaffWechat      string          `json:"check_staff_wechat"`
	WithdrawCopywriting   string          `json:"withdraw_copywriting"`
	RechargeCopywriting   string          `json:"recharge_copywriting"`
	Banners               []Banner        `json:"banners"`
}

//ForceUpdateInfo .
type ForceUpdateInfo struct {
	ForceUpdate bool   `json:"is_force_update"`
	NoticeCount uint64 `json:"notice_count"`
	URL         string `json:"url"`
	Version     string `json:"version"`
	Copywriting string `json:"copywriting"`
}

//ShareInfo .
type ShareInfo struct {
	Rule             string `json:"rule"`
	URL              string `json:"url"`
	AwardCopyWriting string `json:"award_copywriting"`
}

//Gift .
type Gift struct {
	ID       uint64 `json:"gift_id"`
	GiftName string `json:"name"`
	Price    uint64 `json:"price"`
	ImgURL   string `json:"img"`
}

//Product .
type Product struct {
	ID          uint64  `json:"product_id"`
	ProductName string  `json:"name"`
	Extra       uint64  `json:"extra"`
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
	//BannerTypeUserDetail 用户详情
	BannerTypeUserDetail
	//BannerTypeInvite 邀请用户
	BannerTypeInvite
	//BannerTypeRecharge 充值
	BannerTypeRecharge
)

//Banner 首页banner
type Banner struct {
	Image      string `json:"img"`
	URL        string `json:"url"`
	BannerType int    `json:"banner_type"`
}

//TableName .
func (Config) TableName() string {
	return "config"
}

//GetCommonGiftInfo 获取礼物列表
func (c *Config) GetCommonGiftInfo() ([]Gift, error) {
	if tm, ok := updateTime[configTypeGift]; giftList == nil || !ok || tm.Add(time.Minute*5).Before(time.Now()) {
		var cf []Config
		if err := db.Debug().Where("conf_key = ?", configTypeGift).Find(&cf).Error; err != nil {
			return nil, err
		}

		var gf []Gift
		for index := range cf {
			var g Gift
			if err := utils.JSONUnMarshal(cf[index].Value, &g); err != nil {
				beego.Error("礼物信息解析失败", err)
				continue
			}

			g.ID = cf[index].ID
			gf = append(gf, g)
		}

		giftList = gf
		updateTime[configTypeGift] = time.Now()
		return gf, nil
	}
	return giftList, nil
}

//GetProductInfo 获取商品列表
func (c *Config) GetProductInfo() ([]Product, error) {
	if tm, ok := updateTime[configTypeProduct]; giftList == nil || !ok || tm.Add(time.Minute*5).Before(time.Now()) {
		var cf []Config
		if err := db.Debug().Where("conf_key = ?", configTypeProduct).Find(&cf).Error; err != nil {
			return nil, err
		}

		var pf []Product
		for index := range cf {
			var p Product
			if err := utils.JSONUnMarshal(cf[index].Value, &p); err != nil {
				beego.Error("礼物信息解析失败", err)
				continue
			}

			p.ID = cf[index].ID
			pf = append(pf, p)
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

//GetCommonConfig 获取banner
func (c *Config) GetCommonConfig() (*CommonConfig, error) {
	if tm, ok := updateTime[configTypeCommon]; commonConfig == nil || !ok || tm.Add(time.Minute*5).Before(time.Now()) {
		if err := db.Model(c).Where("conf_key = ?", configTypeCommon).First(c).Error; err != nil {
			return nil, err
		}

		var cc CommonConfig
		if err := utils.JSONUnMarshal(c.Value, &cc); err != nil {
			return nil, err
		}

		updateTime[configTypeCommon] = time.Now()
		commonConfig = &cc
		return commonConfig, nil
	}
	return commonConfig, nil
}
