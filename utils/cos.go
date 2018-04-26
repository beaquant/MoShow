package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	qCloud "github.com/MrSong0607/image-go-sdk"

	"github.com/astaxie/beego"
)

var (
	qcloudAppid                   uint
	qcloudSID, qcloudSKey, bucket string = beego.AppConfig.String("qcloudSID"), beego.AppConfig.String("qcloudSKey"), beego.AppConfig.String("qcloudBucket")
	cloud                         *qCloud.PicCloud
	qcToken                       *QCloudToken
)

//QCloudTokenResult .
type QCloudTokenResult struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	CodeDesc string      `json:"codeDesc"`
	Data     QCloudToken `json:"data"`
}

//QCloudToken .
type QCloudToken struct {
	ExpiredTime int64             `json:"expiredTime"`
	Credentials qCloudCredentials `json:"credentials"`
}

type qCloudCredentials struct {
	SessionToken string `json:"sessionToken"`
	TmpSecretID  string `json:"tmpSecretId"`
	TmpSecretKey string `json:"tmpSecretKey"`
}

func init() {
	if qcID, err := beego.AppConfig.Int("qcloudAppID"); err != nil {
		panic(err)
	} else {
		qcloudAppid = uint(qcID)
	}

	cloud = &qCloud.PicCloud{Appid: qcloudAppid, SecretId: qcloudSID, SecretKey: qcloudSKey, Bucket: bucket}
}

//GetTecentImgSignV5 .
func GetTecentImgSignV5() (*QCloudToken, error) {
	if qcToken == nil || qcToken.ExpiredTime <= time.Now().Unix() {
		tk, err := AppTempSession(qcloudSID, qcloudSKey)
		qcToken = &tk.Data
		return qcToken, err
	}
	return qcToken, nil
}

//TransCosToCDN .
func TransCosToCDN(origin string) string {
	return strings.Replace(origin,
		"moshow-1255921343.cos.ap-shanghai.myqcloud.com", "moshow-1255921343.file.myqcloud.com", -1)
}

//ImgPornCheckSingle .
func ImgPornCheckSingle(url string) (bool, error) {
	res, err := cloud.PornDetect(url)
	if err != nil {
		return false, err
	}

	if res.Result == qCloud.PornDetectPorn {
		return true, nil
	}

	return false, nil
}

//AppTempSession .
func AppTempSession(secretID, secretKey string) (*QCloudTokenResult, error) {
	policy := `{"statement":[{"action":["name/cos:*"],"effect":"allow","resource":"*"}],"version":"2.0"}`
	baseURL, err := url.Parse("https://sts.api.qcloud.com/v2/index.php")
	if err != nil {
		return nil, err
	}

	param := make(map[string]string)
	param["name"] = "someone"
	param["policy"] = policy //base64.StdEncoding.EncodeToString([]byte(policy))
	param["durationSeconds"] = "7200"
	param["Timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)

	rand.Seed(time.Now().Unix())
	param["Nonce"] = strconv.Itoa(rand.Int())
	param["SecretId"] = secretID
	param["Action"] = "GetFederationToken"

	sign := genSignature("get", baseURL.Host+baseURL.Path, secretKey, param)
	param["Signature"] = sign
	uParams := url.Values{}

	for k := range param {
		uParams.Add(k, param[k])
	}

	finalURL := baseURL.String() + "?" + uParams.Encode()

	resp, err := http.Get(finalURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tk := &QCloudTokenResult{}
	return tk, JSONUnMarshal(string(bd), tk)
}

func genSignature(method, uri, secretkey string, param map[string]string) string {
	param["SignatureMethod"] = "HmacSHA1"

	var keys []string
	for k := range param {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	str := strings.ToUpper(method) + uri + "?"
	for index := range keys {
		if index > 0 {
			str = str + "&"
		}
		str = str + keys[index] + "=" + param[keys[index]]
	}

	h := hmac.New(sha1.New, []byte(secretkey))
	if _, err := h.Write([]byte(str)); err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
