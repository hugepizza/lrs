package tx

import (
	"encoding/json"

	"github.com/hugepizza/lrs/util"

	"github.com/hugepizza/lrs/img"
	"github.com/hugepizza/lrs/types"
)

const searchAPI = "https://cftweb.3g.qq.com/qqappstore_cgi/searchGame?keyWord=%E7%8B%BC%E4%BA%BA%E6%9D%80&contextData=&pageSize=30"
const lrsAppID = 52565390
const lrsPackageName = "com.tencent.tmgp.yongyong.lr"
const mlPackage = ""

type ret struct {
	Data struct {
		AppList []appData `json:"appList"`
	} `json:"data"`
}

type appData struct {
	AppID   int    `json:"appId"`
	PkgName string `json:"pkgName"`
	LogoURL string `json:"logoUrl"`
	ApkURL  string `json:"apkUrl"`
}

// GetRank .
func GetRank() (*types.RankInfo, error) {
	header := map[string][]string{
		"Referer":          []string{"https://cftweb.3g.qq.com/qqappstore/search"},
		"Accept":           []string{"application/json"},
		"X-Requested-With": []string{"XMLHttpRequest"},
		"User-Agent":       []string{"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1"},
	}
	bs, err := util.HTTPGet(searchAPI, header)
	if err != nil {
		return nil, err
	}
	ret := &ret{}
	err = json.Unmarshal(bs, ret)
	if err != nil {
		return nil, err
	}
	r := &types.RankInfo{}
	for i, v := range ret.Data.AppList {
		if v.PkgName == lrsPackageName || v.AppID == lrsAppID {
			r.LrsRank = i + 1
		}
		same := img.GetSimilar(v.LogoURL)
		if same > 0.8 {
			r.MlRank = i + 1
			r.MlURL = v.ApkURL
		}
		if r.LrsRank > 0 && r.MlRank > 0 {
			break
		}
	}
	return r, nil
}
