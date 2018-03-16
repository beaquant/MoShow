package utils

import (
	sign "boys/common/qcloud/sign"
	"time"

	"github.com/astaxie/beego"
)

var (
	qcloudAppid                   uint
	qcloudSID, qcloudSKey, bucket string = beego.AppConfig.String("qcloudSID"), beego.AppConfig.String("qcloudSKey"), beego.AppConfig.String("qcloudBucket")
)

func init() {
	if qcID, err := beego.AppConfig.Int("qcloudAppID"); err != nil {
		panic(err)
	} else {
		qcloudAppid = uint(qcID)
	}
}

//GetTecentImgSign .
func GetTecentImgSign() (string, error) {
	return sign.AppSignV2(qcloudAppid, qcloudSID, qcloudSKey, bucket, uint(time.Now().Add(time.Hour).Unix()))
}
