package utils

import (
	sign "boys/common/qcloud/sign"
	"time"
)

var (
	qcloudAppid                   uint
	qcloudSID, qcloudSKey, bucket string
)

//GetTecentImgSign .
func GetTecentImgSign() (string, error) {
	return sign.AppSignV2(qcloudAppid, qcloudSID, qcloudSKey, bucket, uint(time.Now().Add(time.Hour).Unix()))
}
