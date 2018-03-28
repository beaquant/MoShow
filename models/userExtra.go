package models

import (
	"MoShow/utils"
	"errors"
	"sort"
	"strconv"

	"github.com/jinzhu/gorm"
)

//UserExtra 用户附加信息
type UserExtra struct {
	ID          uint64 `json:"user_id" gorm:"column:id;primary_key" description:"用户ID"`
	GiftHistory string `json:"-" gorm:"column:gift_his" description:"收到的礼物"`
	IncomeHis   uint64 `json:"income_his" gorm:"column:income_his" description:"历史总收益"`
	BalanceHis  uint64 `json:"balance_his" gorm:"column:balance_his" description:"历史总充值"`
	InviteCount uint64 `json:"invite_count" gorm:"column:invite_count" description:"邀请总人数"`
}

//GiftHisInfo .
type GiftHisInfo struct {
	Count    uint64 `json:"count"`
	GiftInfo Gift   `json:"gift_info"`
}

//GiftHisInfoList .
type GiftHisInfoList []GiftHisInfo

func (g GiftHisInfoList) Len() int           { return len(g) }
func (g GiftHisInfoList) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g GiftHisInfoList) Less(i, j int) bool { return g[i].GiftInfo.Price < g[j].GiftInfo.Price }

//TableName .
func (UserExtra) TableName() string {
	return "user_extra"
}

//Add .
func (u *UserExtra) Add(trans *gorm.DB) error {
	if trans == nil {
		trans = db
	}

	if len(u.GiftHistory) == 0 {
		u.GiftHistory = "{}"
	}

	return trans.Create(u).Error
}

//ReadOrCreate .
func (u *UserExtra) Read() error {
	return db.Find(u).Error
}

//AddGiftCount .
func (u *UserExtra) AddGiftCount(gft Gift, count uint64, trans *gorm.DB) error {
	if trans == nil {
		trans = db
	}

	gstr, err := utils.JSONMarshalToString(&GiftHisInfo{GiftInfo: gft, Count: count})
	if err != nil {
		return err
	}

	idStr := strconv.FormatUint(gft.ID, 10)
	countStr := strconv.FormatUint(count, 10)

	return trans.Model(u).Update("gift_his", gorm.Expr(`if(isnull(gift_his ->>'$."`+idStr+`"'),JSON_SET(COALESCE(gift_his,"{}"),'$."`+idStr+`"',cast(? as json)),JSON_SET(gift_his,'$."`+idStr+`"."Count"',gift_his->>'$."`+idStr+`"."Count"' + `+countStr+`))`, gstr)).Error
}

//GetGiftHis .
func (u *UserExtra) GetGiftHis() ([]GiftHisInfo, error) {
	gftHist := make(map[uint64]GiftHisInfo)
	if err := db.Find(u).Error; err != nil {
		return nil, err
	}

	if len(u.GiftHistory) == 0 {
		return nil, nil
	}

	if err := utils.JSONUnMarshal(u.GiftHistory, &gftHist); err != nil {
		return nil, err
	}

	var gfts []GiftHisInfo
	for _, v := range gftHist {
		gfts = append(gfts, v)
	}

	sort.Sort(GiftHisInfoList(gfts))
	return gfts, nil
}

//AddBalanceHist .
func (u *UserExtra) AddBalanceHist(count uint64, trans *gorm.DB) error {
	if u.ID == 0 {
		return errors.New("user_extra 更新用户历史余额 必须指定用户ID")
	}

	if trans != nil {
		return trans.Model(u).Update("balance_his", gorm.Expr("balance_his + ?", count)).Error
	}
	return db.Model(u).Update("balance_his", gorm.Expr("balance_his + ?", count)).Error
}

//AddIncomeHist .
func (u *UserExtra) AddIncomeHist(count uint64, trans *gorm.DB) error {
	if u.ID == 0 {
		return errors.New("user_extra 更新用户历史收益 必须指定用户ID")
	}

	if trans != nil {
		return trans.Model(u).Update("income_his", gorm.Expr("income_his + ?", count)).Error
	}
	return db.Model(u).Update("income_his", gorm.Expr("income_his + ?", count)).Error
}

//AddInviteCount .
func (u *UserExtra) AddInviteCount(trans *gorm.DB) error {
	if u.ID == 0 {
		return errors.New("user_extra 更新邀请人数 必须指定用户ID")
	}

	if trans != nil {
		return trans.Model(u).Update("invite_count", gorm.Expr("invite_count + ?", 1)).Error
	}
	return db.Model(u).Update("invite_count", gorm.Expr("invite_count + ?", 1)).Error
}
