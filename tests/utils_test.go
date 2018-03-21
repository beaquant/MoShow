package test

import (
	"MoShow/models"
	"MoShow/utils"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/silenceper/wechat/oauth"
)

func TestMsgSend(t *testing.T) {
	res, err := utils.SendMsgByAPIKey("18868875634", "短信测试")
	t.Log(res, err)
}

func TestPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
		}
	}()

	panic("test")
}

func TestRegex(t *testing.T) {
	r, _ := regexp.Compile("/v1/.+?/")
	ss := r.FindStringSubmatch("/v1/auth/18868875634/sendcode")
	if ss == nil || len(ss) == 0 {
		panic(ss)
	}
	panic(ss[0])
}

func TestArrayJoin(t *testing.T) {
	array := []string{"1", "2", "3"}
	t.Log(strings.Join(array, "','"))
}

func TestURLParse(t *testing.T) {
	if _, err := url.ParseRequestURI(""); err != nil {
		t.Error(err)
	}
}

func TestCompute(t *testing.T) {
	var a uint64
	a = 5
	t.Log(a * 3 / 10)
}

func TestJson(t *testing.T) {
	g := &models.Product{}
	t.Log(utils.JSONMarshalToString(g))
}

func TestSlice(t *testing.T) {
	a := []string{""}
	b := a[:0]

	str, err := utils.JSONMarshalToString(&b)
	t.Log(str, err)
}

func TestContains(t *testing.T) {
	t.Log(strings.Contains("blueMr/1.1.0 (iPhone; iOS 11.2.5; Scale/2.00)", "iPhone"))
}

func TestPointAddr(t *testing.T) {
	a := &struct{}{}
	b := &struct{}{}
	t.Log(a == b)
}

func TestWechatLogin(t *testing.T) {
	o := oauth.NewOauth(nil)
	info, err := o.GetUserInfo("7_9P0JWuyX3Hq6iQfDzfZEISor6ErwjfKD7Hz61sErhK819sP7-j4oe30881axlSORBSDX_XU98K4oepcvgaZK05ZLCT1XCae5pts_tNk8LbU", "oaGlO1gTBzsnPbDWEyZgZdrq17Do")
	t.Log(info, err)
}

func TestStringReplace(t *testing.T) {
	str := strings.Replace("http://bluemr-1254204939.cossh.myqcloud.co/photo/44919201802230116495945-1000x1777.jpg",
		"bluemr-1254204939.cossh.myqcloud.com", "bluemr-1254204939.file.myqcloud.com", -1)
	t.Log(str)
}

func TestCosSign(t *testing.T) {

}

func TestPornImg(t *testing.T) {
	utils.ImgPornCheckSingle("http://bluemr-1254204939.cossh.myqcloud.com/img/1521621081565194700.jpg")
}
