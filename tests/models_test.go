package test

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"MoShow/models"
	"testing"
)

func TestUsers(t *testing.T) {
	u := &models.User{ID: 1}
	err := u.ReadFromPhoneNumber()
	if err != nil {
		panic(err)
	}
}