package test

import (
	"MoShow/utils"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"testing"
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
