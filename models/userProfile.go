package models

import (
	"MoShow/utils"
	"errors"
	"strconv"
	"time"

	"github.com/astaxie/beego"

	"github.com/jinzhu/gorm"
)

const (
	//GenderDefault .
	GenderDefault = iota
	//GenderMan .
	GenderMan
	//GenderWoman .
	GenderWoman
)

const (
	//UserTypeNormal 普通用户
	UserTypeNormal = iota
	//UserTypeAnchor 主播
	UserTypeAnchor
)

const (
	//UserStatusNormal 正常
	UserStatusNormal = iota
	//UserStatusHot 推荐
	UserStatusHot
	//UserStatusBlock 屏蔽
	UserStatusBlock
)

//UserProfile .
type UserProfile struct {
	ID          uint64 `json:"user_id" gorm:"column:id;primary_key" description:"用户ID"`
	Alias       string `json:"alias" gorm:"column:alias" description:"昵称"`
	Gender      int    `json:"gender" gorm:"column:gender" description:"性别"`
	CoverPic    string `json:"-" gorm:"column:cover" description:"形象展示,包括头像,相册,视频"`
	Description string `json:"description" gorm:"column:description" description:"签名"`
	Birthday    int64  `json:"birthday" gorm:"column:birthday" description:"生日"`
	Location    string `json:"location" gorm:"column:location" description:"地区"`
	Balance     uint64 `json:"balance" gorm:"column:balance" description:"余额"`
	Income      uint64 `json:"income" gorm:"column:income" description:"收益"`
	Price       uint64 `json:"price" gorm:"column:price" description:"视频价格/分"`
	UserType    int    `json:"user_type" gorm:"column:user_type" description:"用户类型"`
	UserStatus  int    `json:"-" gorm:"column:user_status" description:"用户状态"`
	UpdateAt    int64  `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	ImToken     string `json:"-" gorm:"column:im_token" description:"网易云信token"`
	Followers   string `json:"-" gorm:"column:follower" description:"关注者"`
	Following   string `json:"-" gorm:"column:following" description:"正在关注"`
}

//UserCoverInfo .
type UserCoverInfo struct {
	CoverPicture *Picture  `json:"cover_pic_info"`
	DesVideo     *Video    `json:"video_info"`
	Gallery      []Picture `json:"gallery"`
}

//Picture .
type Picture struct {
	ImageURL string `json:"image_url"`
	Disable  bool   `json:"disable"`
	Checked  bool   `json:"checked"`
}

//Video .
type Video struct {
	VideoURL string `json:"video_url"`
	Disable  bool   `json:"disable"`
	Checked  bool   `json:"checked"`
}

//FollowInfo .
type FollowInfo struct {
	FollowTime int64
}

//TableName .
func (UserProfile) TableName() string {
	return "user_profile"
}

//ToString .
func (u *UserCoverInfo) ToString() string {
	str, err := utils.JSONMarshalToString(u)
	if err != nil {
		beego.Error(err)
		return ""
	}
	return str
}

//Add .
func (u *UserProfile) Add(trans *gorm.DB) error {
	if trans != nil {
		return trans.Create(u).Error
	}
	return db.Create(u).Error
}

func (u *UserProfile) Read() error {
	if u.ID == 0 {
		return errors.New("必须指定user_profile的id")
	}
	return db.Where("id = ?", u.ID).Find(u).Error
}

//Update .
func (u *UserProfile) Update(fields map[string]interface{}) error {
	fields["update_at"] = time.Now().Unix()
	return db.Model(u).Updates(fields).Error
}

//AddFollow 添加关注
func (u *UserProfile) AddFollow(id uint64) error {
	idStr := strconv.FormatUint(id, 10)
	fis, _ := utils.JSONMarshalToString(&FollowInfo{FollowTime: time.Now().Unix()})

	trans := db.Begin()
	if err := trans.Model(u).Updates(map[string]interface{}{"following": `JSON_SET(COALESCE(following,'{}'),'$."` + idStr + `"',CAST('` + fis + `' AS JSON))`}).Error; err != nil {
		trans.Rollback()
		return err
	}

	if err := trans.Model(&UserProfile{ID: id}).Updates(map[string]interface{}{"follower": `JSON_SET(COALESCE(follower,'{}'),'$."` + idStr + `"',CAST('` + fis + `' AS JSON))) `}).Error; err != nil {
		trans.Rollback()
		return err
	}

	trans.Commit()
	return nil
}

//UnFollow 取消关注
func (u *UserProfile) UnFollow(id uint64) error {
	idStr := strconv.FormatUint(id, 10)
	trans := db.Begin()
	if err := trans.Model(u).Updates(map[string]interface{}{"following": `JSON_REMOVE(follower,'$."` + idStr + `"')`}).Error; err != nil {
		trans.Rollback()
		return err
	}

	if err := trans.Model(&UserProfile{ID: id}).Updates(map[string]interface{}{"follower": `JSON_REMOVE(follower,'$."` + idStr + `"')`}).Error; err != nil {
		trans.Rollback()
		return err
	}

	trans.Commit()
	return nil
}

//GetCover .
func (u *UserProfile) GetCover() *UserCoverInfo {
	if len(u.CoverPic) > 0 {
		ucp := &UserCoverInfo{}
		if err := utils.JSONUnMarshal(u.CoverPic, ucp); err != nil {
			beego.Error(err)
			return nil
		}
		return ucp
	}
	return nil
}

//GetFollowers .
func (u *UserProfile) GetFollowers() map[uint64]FollowInfo {
	if len(u.Followers) > 0 {
		fl := make(map[uint64]FollowInfo)
		if err := utils.JSONUnMarshal(u.Followers, &fl); err != nil {
			beego.Error(err)
			return nil
		}
		return fl
	}
	return nil
}

//GetFollowing .
func (u *UserProfile) GetFollowing() map[uint64]FollowInfo {
	if len(u.Following) > 0 {
		fl := make(map[uint64]FollowInfo)
		if err := utils.JSONUnMarshal(u.Following, &fl); err != nil {
			beego.Error(err)
			return nil
		}
		return fl
	}
	return nil
}

//GetInviteList 获取我邀请的用户列表
func (u *UserProfile) GetInviteList() (ul []UserProfile, err error) {
	if err = db.Joins("left join users on user_profile.id = users.id").Where("invited_by = ?", u.ID).Find(&ul).Error; err != nil {
		return nil, err
	}

	return
}

//AddBalance .
func (u *UserProfile) AddBalance(amount int, trans *gorm.DB) error {
	if trans != nil {
		return trans.Model(u).Update("balance", gorm.Expr("balance + ?", amount)).Error
	}
	return db.Model(u).Update("balance", gorm.Expr("balance + ?", amount)).Error
}

//AddIncome .
func (u *UserProfile) AddIncome(amount int, trans *gorm.DB) error {
	if trans != nil {
		return trans.Model(u).Update("income", gorm.Expr("income + ?", amount)).Error
	}
	return db.Model(u).Update("income", gorm.Expr("income + ?", amount)).Error
}

//AllocateFund 划款
func (u *UserProfile) AllocateFund(to, invitation *UserProfile, amount, incomeAmount, inviteAmount uint64, trans *gorm.DB) error {
	if u.Balance+u.Income < amount { //检查余额
		return errors.New("用户余额不足，扣款(" + strconv.FormatUint(amount, 10) + ")失败,所有钱包余额合计:" + strconv.FormatUint(u.Balance+u.Income, 10))
	}

	if u.Balance < amount { //余额钱包金额足够扣款
		if err := u.AddBalance(-int(amount), trans); err != nil { //扣款
			return errors.New("发起人扣款失败\t" + err.Error())
		}
	} else { //余额钱包金额不足以扣款
		if err := u.AddBalance(-int(u.Balance), trans); err != nil { //从余额钱包扣款
			return errors.New("发起人扣款失败\t" + err.Error())
		}

		if err := u.AddIncome(-int(amount-u.Balance), trans); err != nil { //从收益钱包扣款
			return errors.New("发起人扣款失败\t" + err.Error())
		}
	}

	if err := to.AddIncome(int(incomeAmount), trans); err != nil { //增加余额
		return errors.New("接受人增加余额失败\t" + err.Error())
	}

	if invitation != nil {
		if uint64(incomeAmount*3/10) < inviteAmount {
			return errors.New("收益分成比例异常,收入:" + strconv.FormatUint(incomeAmount, 10) + ",分成:" + strconv.FormatUint(inviteAmount, 10))
		}

		if err := invitation.AddIncome(int(inviteAmount), trans); err != nil { //增加余额
			return errors.New("邀请人增加余额失败\t" + err.Error())
		}
	}

	return nil
}
