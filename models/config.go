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
	configTypeCkModeRegs = "checkMode_regs"
)

const (
	//BannerTypeImg 纯图片banner
	BannerTypeImg = iota
	//BannerTypeLink 内链
	BannerTypeLink
	//BannerTypeUserDetail 用户详情
	BannerTypeUserDetail
	//BannerTypeInvite 邀请用户
	BannerTypeInvite
	//BannerTypeRecharge 充值
	BannerTypeRecharge
	//BannerTypeHTTPLink 外链
	BannerTypeHTTPLink
)

var (
	giftList      []Gift
	productList   []Product
	checkModeRegs []string
	updateTime    = make(map[string]time.Time)
	incomeRate    *IncomeRate
	commonConfig  *CommonConfig
)

//Config .
type Config struct {
	ID    uint64 `json:"id" gorm:"column:id;primary_key"`
	Key   string `json:"key" gorm:"column:conf_key"`
	Value string `json:"val" gorm:"column:val"`
}

//CommonConfig 通用配置
type CommonConfig struct {
	AnchorVideoRecord     bool             `json:"ac_video_record" description:"是否开启主播录制视频"`     //是否开启主播录制视频
	UserVideoRecordbool   bool             `json:"user_video_record" description:"是否开启用户录制视频"`   //是否开启用户录制视频
	UserProtocol          string           `json:"user_protocal" description:"用户协议"`             //用户协议
	ForceUpdate           *ForceUpdateInfo `json:"force_update,omitempty" description:"强制更新"`    //强制更新
	Share                 ShareInfo        `json:"share" description:"邀请"`                       //邀请
	CustomerServiceWechat string           `json:"customer_service_wechat" description:"客服人员微信"` //客服人员微信
	CheckStaffWechat      string           `json:"check_staff_wechat" description:"审核人员微信"`      //审核人员微信
	WithdrawCopywriting   string           `json:"wd_copywriting" description:"提现文案"`            //提现文案
	RechargeCopywriting   string           `json:"rcg_copywriting" description:"充值文案"`           //充值文案
	Banners               []Banner         `json:"banners" description:"轮播图"`                    //轮播图
	VideoPrice            uint64           `json:"vod_value" description:"形象视频扣费价格"`             //形象视频扣费价格
	MessagePrice          uint64           `json:"msg_value" description:"私聊扣费价格"`               //私聊扣费价格
}

//ForceUpdateInfo .
type ForceUpdateInfo struct {
	ForceUpdate bool   `json:"is_force_update" description:"是否强制更新"`
	NoticeCount uint64 `json:"notice_count" description:"强制更新提醒次数"`
	URL         string `json:"url" description:"下载地址"`
	Version     string `json:"version" description:"版本号"`
	Copywriting string `json:"copywriting" description:"文案"`
}

//ShareInfo .
type ShareInfo struct {
	Rule string `json:"rule" description:"奖励规则"`
	// URL              string   `json:"url" description:"链接"`
	AwardCopyWriting []string `json:"award_copywriting" description:"奖励文案"`
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

//Banner 首页banner
type Banner struct {
	Image      string `json:"img"`
	UserID     uint64 `json:"user_id,omitempty"`
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
			g.ImgURL = utils.TransCosToCDN(g.ImgURL)
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

//GetcheckModeRegs 获取审核模式正则串
func (c *Config) GetcheckModeRegs() ([]string, error) {
	if tm, ok := updateTime[configTypeCkModeRegs]; checkModeRegs == nil || !ok || tm.Add(time.Minute*5).Before(time.Now()) {
		if err := db.Model(c).Where("conf_key = ?", configTypeCkModeRegs).First(c).Error; err != nil {
			return nil, err
		}

		var cc []string
		if err := utils.JSONUnMarshal(c.Value, &cc); err != nil {
			return nil, err
		}

		updateTime[configTypeCkModeRegs] = time.Now()
		checkModeRegs = cc
		return checkModeRegs, nil
	}
	return checkModeRegs, nil
}
