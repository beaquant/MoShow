package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//IPApiPoint .
const IPApiPoint = "http://ip-api.com/json/%s"

//IPInfo .
type IPInfo struct {
	Address     string `json:"as"`
	City        string `json:"city"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Isp         string `json:"isp"`
	Org         string `json:"org"`
	Query       string `json:"query"`
	Region      string `json:"region"`
	RegionName  string `json:"regionName"`
	Status      string `json:"status"`
	Timezone    string `json:"timezone"`
}

//GetIPInfo .
func GetIPInfo(ip string) (*IPInfo, error) {
	resp, err := http.Get(fmt.Sprintf(IPApiPoint, ip))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ii := &IPInfo{}
	if err := JSONUnMarshalFromByte(bd, ii); err != nil {
		return nil, err
	}
	return ii, nil
}
