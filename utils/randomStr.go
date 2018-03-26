package utils

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/satori/go.uuid"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

//RandStringBytesMaskImprSrc .
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// RandNumber .
func RandNumber(min, max int) int {
	return min + rand.Intn(max-min)
}

//UUIDBase64String .
func UUIDBase64String() (string, error) {
	u1, err := uuid.NewV4()
	return base64.URLEncoding.EncodeToString(u1.Bytes()), err
}
