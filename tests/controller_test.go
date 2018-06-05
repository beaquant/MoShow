package test

import (
	"os"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"MoShow/controllers"
	"MoShow/utils"
	"testing"
)

func TestCheckModePattern(t *testing.T) {
	os.Setenv("GOCACHE", "off")
	t.Log(controllers.IsCheckMode("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36 QIHU 360SE/Nutch-1.13"))
}

func TestActive(t *testing.T) {
	os.Setenv("GOCACHE", "off")
	go func() { t.Log(controllers.SendActivity(169298)) }()
	go func() { t.Log(controllers.SendActivity(169293)) }()
	go func() { t.Log(controllers.SendActivity(169143)) }()
	go func() { t.Log(controllers.SendActivity(170034)) }()
	time.Sleep(150 * time.Second)
}
func TestToken(t *testing.T) {
	tk := &controllers.Token{}
	t.Log(tk.Decrypt("iA8aheaDexLkxoijX4Yvz0vq_C52QcSmVCc6BlUpy8WpcDDZdG8s_XpO4pvAyJcIsDrxhjn7mWNJRrQJTGc0CGCMdFsWW9OiDQZAVGL7x40="))
	t.Log(utils.JSONMarshalToString(tk))
}
