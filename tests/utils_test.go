package test

import (
	"MoShow/models"
	"MoShow/utils"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	sign "github.com/MrSong0607/image-go-sdk/sign"
	"github.com/silenceper/wechat/oauth"
	"github.com/sirupsen/logrus"
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

	// a := &struct{}{}
	// t.Log(utils.JSONUnMarshal("null", a))

	// a := &models.Banner{}
	// a := &models.UserCoverInfo{CoverPicture: &models.Picture{ImageURL: "1"}, DesVideo: &models.Video{VideoURL: "1"}}
	a := &models.ActiveDetail{}
	t.Log(utils.JSONMarshalToString(a))
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
	s, err := sign.AppSignV5("AKIDJYdWDvGCiCdxFUwqtxyujZXVPI9ztbDe", "RlxUvFPOCoIguzugHoBNhsshP1To1X6a", "/photo", "put", 3600*24)
	t.Log(s, err)
}

func TestPornImg(t *testing.T) {
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:8889")
	t.Log(utils.ImgPornCheckSingle("http://bluemr-1254204939.cossh.myqcloud.com/photo/117516201804121657141391-959x1317.jpg"))
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

func TestImImgMsg(t *testing.T) {
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:8889")
	t.Log(utils.SendP2PSysImageMessage("http://pic6.wed114.cn/20130404/20130404143859739.JPG",
		[]string{"169138", "169136", "169137", "169134", "100008", "5"}))
}

func TestImvoiceMsg(t *testing.T) {
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:8889")
	t.Log(utils.SendP2PSysVoiceMessage("https://moshow-1255921343.cos.ap-shanghai.myqcloud.com/voice/f272948ce9d076f4f0f5bd83872af5e2.aac", 10*1000,
		[]string{"169178"}))
}

func TestImVideoMsg(t *testing.T) {
	// os.Setenv("HTTP_PROXY", "http://127.0.0.1:8889")
	t.Log(utils.SendP2PSysVideoMessage("https://moshow-1255921343.file.myqcloud.com/video/1691361525878287277977386.mp4", 10,
		[]string{"169298"}))
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

func TestOsFileMode(t *testing.T) {
	t.Logf("ModePerm:%o", os.ModePerm)
	t.Logf("ModeAppend:%o", os.ModeAppend)
	t.Logf("ModeCharDevice:%o", os.ModeCharDevice)
	t.Logf("ModeDevice:%o", os.ModeDevice)
	t.Logf("ModeDir:%o", os.ModeDir)
	t.Logf("ModeExclusive:%o", os.ModeExclusive)
	t.Logf("ModeNamedPipe:%o", os.ModeNamedPipe)
	t.Logf("ModeSetgid:%o", os.ModeSetgid)
	t.Logf("ModeSetuid:%o", os.ModeSetuid)
	t.Logf("ModeSocket:%o", os.ModeSocket)
	t.Logf("ModeSticky:%o", os.ModeSticky)
	t.Logf("ModeSymlink:%o", os.ModeSymlink)
	t.Logf("ModeTemporary:%o", os.ModeTemporary)
	t.Logf("ModeType:%o", os.ModeType)
}

func TestGetIP(t *testing.T) {
	if i, err := utils.GetIPInfo("45.78.22.174"); err != nil {
		t.Error(err)
	} else {
		t.Log(utils.JSONMarshalToString(i))
	}
}

func TestLogrus(t *testing.T) {
	log := logrus.WithFields(logrus.Fields{"dial_id": 123, "ext": "asdad"})
	logrus.SetFormatter(TextFormatter{})
	log.Error("test")
	log.Warn("test2")
	log.Info("ss")
}

type TextFormatter struct {
}

func (TextFormatter) Format(e *logrus.Entry) ([]byte, error) {
	str := fmt.Sprintf("%s[%s] [%d] %s", e.Time.Format("06/01/02 15:04:05"), strings.ToUpper(string(e.Level.String()[0])), e.Data["dial_id"], e.Message)
	for k, v := range e.Data {
		if k != "dial_id" {
			str = fmt.Sprintf("%s %s:%s", str, k, v)
		}
	}
	str += `
`
	return []byte(str), nil
}

func Test素数(t *testing.T) {
	a := 23

	if a < 1 {
		t.Log("输入错误,不能小于1")
		return
	}

	if a < 3 {
		t.Log("是素数")
		return
	}

	for i := 2; i < a; i++ {
		if a%i == 0 {
			t.Log("不是素数")
			return
		}
	}

	t.Log("是素数")
}

func TestUint(t *testing.T) {
	a := uint64(10)
	t.Log(a - 15)
}
