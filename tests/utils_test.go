package test

import (
	"MoShow/utils"
	"fmt"
	"net/url"
	"regexp"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/silenceper/wechat/oauth"
)

func TestMsgSend(t *testing.T) {
	res, err := utils.SendMsgByAPIKey("18868875634", "短信测试")
	t.Log(res, err)
}

func TestPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			debug.PrintStack()
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
	// g := make(map[uint64]models.GiftHisInfo)
	// g[2] = models.GiftHisInfo{Count: 10, GiftInfo: models.Gift{ID: 2}}
	// t.Log(utils.JSONMarshalToString(g))

	a := &struct{}{}
	t.Log(utils.JSONUnMarshal("null", a))
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

func TestReg(t *testing.T) {
	var re = regexp.MustCompile("\\d+")
	t.Log(re.FindString("BlueMr1561"))
}

func TestImTextMesg(t *testing.T) {
	t.Log(utils.SendP2PMessage("1", "3", "normal message test"))
}

func TestImSysMsg(t *testing.T) {
	// t.Log(utils.SendSysMessage(&utils.ImSysNotifyMessage{Message: "名门97", Type: utils.ImSysNotifyMessageTypeText}, []string{"4"}))
}

func TestTimeSubSeconds(t *testing.T) {
	start := time.Now().Unix()
	time.Sleep(time.Second * 2)

	stop := time.Now()
	tl := stop.Sub(time.Unix(start, 0)).Seconds()
	t.Log(tl)
	t.Log(uint64(tl))

	itl := stop.Unix() - start
	t.Log("unix time sub", itl)
}

func TestMapDelete(t *testing.T) {
	dic := make(map[uint64]interface{})

	val, ok := dic[1]
	t.Log(val, ok)
}

func TestUint64Sub(t *testing.T) {
	a := uint64(10)
	t.Log(-time.Duration(a) * time.Second)
}
