package test

import (
	"os"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"MoShow/controllers"
	"testing"
)

func TestCheckModePattern(t *testing.T) {
	os.Setenv("GOCACHE", "off")
	t.Log(controllers.IsCheckMode("Mxiu/1.0.0(Android:24;VTR-AL00)xiaomi"))
}

func TestActive(t *testing.T) {
	os.Setenv("GOCACHE", "off")
	go func() { t.Log(controllers.SendActivity(169298)) }()
	go func() { t.Log(controllers.SendActivity(169293)) }()
	go func() { t.Log(controllers.SendActivity(169143)) }()
	go func() { t.Log(controllers.SendActivity(170034)) }()
	time.Sleep(50 * time.Second)
}
