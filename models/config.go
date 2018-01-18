package models

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
