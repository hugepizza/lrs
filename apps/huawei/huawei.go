package huawei

import (
	"encoding/json"

	"github.com/hugepizza/lrs/img"
	"github.com/hugepizza/lrs/types"
	"github.com/hugepizza/lrs/util"
)

const searchAPI = "https://appgallery.cloud.huawei.com/uowap/index?method=internal.getTabDetail&maxResults=25&reqPageNum=1&serviceType=13&uri=searchApp%7C%E7%8B%BC%E4%BA%BA%E6%9D%80"
const lrsAPPID = "C100104305"
const mlAPPID = "C10740813"

type ret struct {
	LayoutData []struct {
		DataList   []appData `json:"dataList"`
		LayoutName string    `json:"layoutName"`
	} `json:"layoutData"`
}

type appData struct {
	AppID   string `json:"appid"`
	DownURL string `json:"downurl"`
	Icon    string `json:"icon"`
}

// GetRank .
func GetRank() (*types.RankInfo, error) {
	header := map[string][]string{
		"Host":       []string{"appgallery.cloud.huawei.com"},
		"Accept":     []string{"application/json"},
		"User-Agent": []string{"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1"},
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
	ranks := make([]appData, 0, 10)
	for _, v := range ret.LayoutData {
		if v.LayoutName == "safeappcard" {
			ranks = append(ranks, v.DataList...)
		}
	}
	r := &types.RankInfo{}
	for i, v := range ranks {
		if v.AppID == lrsAPPID {
			r.LrsRank = i + 1
		}
		same := img.GetSimilar(v.Icon)
		if same > 0.8 {
			r.MlRank = i + 1
			r.MlURL = v.DownURL
		}
		if r.LrsRank > 0 && r.MlRank > 0 {
			break
		}
	}
	return r, nil
}
