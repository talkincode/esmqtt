package cziploc

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/guonaihong/gout"
	"github.com/talkincode/esmqtt/common"
)

const fetchUrl = "http://cz88.rtbasia.com/getDownloadInfo"

/*
	{
	  "code": 200,
	  "success": true,
	  "message": "操作成功",
	  "data": {
	    "version": "v20220727",
	    "resDownloadTime": 3,
	    "downloadUrl": "https://alpha.cz88.net/api/datDownload/download?downloadCode=eb0acb9a-1b67-3882-a49c-*****",
	    "maxDownloadTime": 3
	  },
	  "time": "2022-07-27 17:15:34"
	}
*/
type downloadInfo struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Version         string `json:"version"`
		ResDownloadTime int    `json:"resDownloadTime"`
		DownloadUrl     string `json:"downloadUrl"`
		MaxDownloadTime int    `json:"maxDownloadTime"`
	}
}

func FetchDat(savepath string) error {
	datinfo, err := os.Stat(savepath)
	if err == nil && time.Now().Sub(datinfo.ModTime()).Hours() < (24*3) {
		fmt.Println("dat file is not expired")
		return nil
	}

	resp := new(downloadInfo)
	err = gout.
		GET(fetchUrl).
		Debug(Debug).
		SetHeader(gout.H{"Authorization": "APPCODE " + Appcode}).
		SetTimeout(time.Second * 5).
		BindJSON(&resp).
		Do()

	if err != nil {
		return err
	}
	if resp.Code != 200 {
		return fmt.Errorf("fetch downloadInfo error %s", resp.Message)
	}

	if resp.Data.ResDownloadTime <= 0 {
		return fmt.Errorf("fetch ResDownloadTime limit")
	}

	body := make([]byte, 0)
	err = gout.GET(resp.Data.DownloadUrl).
		SetTimeout(time.Second * 600).
		BindBody(&body).
		Do()
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return fmt.Errorf("download dat file error")
	}
	tmpfile := path.Join(os.TempDir(), "cziploc.dat")
	err = os.WriteFile(tmpfile, body, 0644)
	if err != nil {
		return err
	}
	if common.FileExists(savepath) {
		_ = common.Copy(savepath, savepath+".bak")
	}
	err = common.Copy(tmpfile, savepath)
	if err == nil {
		return os.Remove(tmpfile)
	}
	return nil
}
