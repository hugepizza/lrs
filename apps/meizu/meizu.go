package meizu

import (
	"encoding/json"

	"github.com/hugepizza/lrs/types"
	"github.com/hugepizza/lrs/util"
)

const searchAPI = "http://app.meizu.com/apps/public/search/page?cat_id=1&keyword=%E7%8B%BC%E4%BA%BA%E6%9D%80&start=0&max=18"
const lrsPackage = ""
const mlPackage = "com.c2vl.kgamebox"

type ret struct {
	Value struct {
		List []appData `json:"list"`
	} `json:"value"`
}

type appData struct {
	PackageName string `json:"package_name"`
	Icon        string `json:"icon"`
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
	for i, v := range ret.Value.List {
		if v.PackageName == mlPackage {
			r.MlRank = i + 1
		}
	}
	return r, nil
}
