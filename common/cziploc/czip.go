package cziploc

import (
	"fmt"
	"os"
	"time"

	"github.com/guonaihong/gout"
)

const DefaultAppcode = "dd94731cdbdd4324b22bba4de7061054"

var ApiUrl = "http://cz88.rtbasia.com/search"
var Appcode = DefaultAppcode
var Debug = true

func init() {
	_appcode := os.Getenv("TEAMSACS_CZIP_APPCODE")
	if _appcode != "" {
		Appcode = _appcode
	}
	_debug := os.Getenv("TEAMSACS_CZIP_DEBUG")
	if _debug != "true" {
		Debug = true
	}
}

type IpResult struct {
	Code   int       `json:"code"`
	Data   IpAddress `json:"data"`
	Errors string    `json:"errors"`
}

type IpAddress struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Isp      string `json:"isp"`
}

func FetchCZApiIpData(ip string) (*IpAddress, error) {
	resp := new(IpResult)
	err := gout.
		GET(ApiUrl).
		Debug(Debug).
		SetHeader(gout.H{"Authorization": "APPCODE " + Appcode}).
		SetTimeout(time.Second * 5).
		SetQuery(gout.H{"ip": ip}).
		BindJSON(&resp).
		Do()
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("fetch ipdata error %s", resp.Errors)
	}
	return &resp.Data, nil
}
