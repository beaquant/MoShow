package models

const (
	//GenderMan .
	GenderMan = iota
	//GenderWoman .
	GenderWoman
)

//UserProfile .
type UserProfile struct {
	ID     uint64 `json:"user_id" gorm:"column:id;primary_key"`
	Alias  string `json:"alias" gorm:"column:alias"`
	Gender int    `json:"gender" gorm:"column:gender"`
}

//TableName .
func (UserProfile) TableName() string {
	return "users"
}

//Add .
func (u *UserProfile) Add() error {
	return db.Create(u).Error
}
