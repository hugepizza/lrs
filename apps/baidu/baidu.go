package baidu

import (
	"encoding/json"

	"github.com/hugepizza/lrs/img"
	"github.com/hugepizza/lrs/types"
	"github.com/hugepizza/lrs/util"
)

const searchAPI = "https://mobile.baidu.com/api/appsearch?word=%E7%8B%BC%E4%BA%BA%E6%9D%80"
const lrsPackage = "com.netease.lrs.baidu"
const mlPackage = ""

type ret struct {
	Data struct {
		Data []appData `json:"data"`
	} `json:"data"`
}

type appData struct {
	Package     string `json:"package"`
	Icon        string `json:"icon"`
	DownloadURL string `json:"download_url"`
}

// GetRank .
func GetRank() (*types.RankInfo, error) {
	bs, err := util.HTTPGet(searchAPI, nil)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	ret := &ret{}
	err = json.Unmarshal(bs, ret)
	if err != nil {
		return nil, err
	}
	r := &types.RankInfo{}
	for i, v := range ret.Data.Data {
		if v.Package == lrsPackage {
			r.LrsRank = i + 1
		}
		same := img.GetSimilar("http:" + v.Icon)
		if same > 0.8 {
			r.MlRank = i + 1
			r.MlURL = v.DownloadURL
		}
		if r.LrsRank > 0 && r.MlRank > 0 {
			break
		}
	}
	return r, nil
}
