package models

const (
	//DialStatusFail 通话状态,失败
	DialStatusFail = iota
	//DialStatusSuccess 成功
	DialStatusSuccess
	//DialStatusException 异常
	DialStatusException
)

//Dial .
type Dial struct {
	ID         uint64 `json:"id" gorm:"column:id;primary_key"`
	FromUserID uint64 `json:"from_user_id" gorm:"column:from_user_id"`
	ToUserID   uint64 `json:"to_user_id" gorm:"column:to_user_id"`
	Duration   int    `json:"duration" gorm:"column:duration"`
	CreateAt   int64  `json:"create_at" gorm:"column:create_at"`
	Status     int    `json:"success" gorm:"column:success"`
	Tag        string `json:"tag" gorm:"column:tag"`
}

//TableName .
func (Dial) TableName() string {
	return "dial"
}

//Add .
func (d *Dial) Add() error {
	return db.Model(d).Create(d).Error
}

//Read .
func (d *Dial) Read() error {
	return db.Where("id = ?", d.ID).Find(d).Error
}

//GetDialList .
func (d *Dial) GetDialList(uid uint64, limit, skip int) ([]Dial, error) {
	if limit == 0 {
		limit = 20
	}

	var lst []Dial
	return lst, db.Where("from_user_id = ?", uid).Find(&lst).Order("create_at desc").Limit(limit).Offset(skip).Error
}

//Del .
func (d *Dial) Del() error {
	return db.Delete(d).Error
}
