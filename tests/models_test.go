package test

import (
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
