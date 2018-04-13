package utils

import (
	"time"

	qCloud "github.com/MrSong0607/image-go-sdk"
	sign "github.com/MrSong0607/image-go-sdk/sign"

	"github.com/astaxie/beego"
)

var (
	qcloudAppid                   uint
	qcloudSID, qcloudSKey, bucket string = beego.AppConfig.String("qcloudSID"), beego.AppConfig.String("qcloudSKey"), beego.AppConfig.String("qcloudBucket")
	cloud                         *qCloud.PicCloud
)

func init() {
	if qcID, err := beego.AppConfig.Int("qcloudAppID"); err != nil {
		panic(err)
	} else {
		qcloudAppid = uint(qcID)
	}

	cloud = &qCloud.PicCloud{Appid: qcloudAppid, SecretId: qcloudSID, SecretKey: qcloudSKey, Bucket: bucket}
}

//GetTecentImgSign .
func GetTecentImgSign() (string, error) {
	return sign.AppSignV2(qcloudAppid, qcloudSID, qcloudSKey, bucket, uint(time.Now().Add(time.Hour).Unix()))
}

//GetTecentImgSignV5 .
func GetTecentImgSignV5(dir string) (string, error) {
	return sign.AppSignV5(qcloudSID, qcloudSKey, dir, "put", 3600*24)
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
