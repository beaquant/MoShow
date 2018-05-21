package test

import (
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"MoShow/controllers"
	"testing"
)

func TestCheckModePattern(t *testing.T) {
	os.Setenv("GOCACHE", "off")
	t.Log(controllers.IsCheckMode("Mxiu/1.0.0(Android:24;VTR-AL00)xiaomi"))
}
