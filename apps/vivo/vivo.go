package vivo

import (
	"encoding/json"

	"github.com/hugepizza/lrs/types"
	"github.com/hugepizza/lrs/util"
)

const searchAPI = "http://game.vivo.com.cn/api/searchGame?pageIndex=1&word=%E7%8B%BC%E4%BA%BA%E6%9D%80&platformVersion=&reqChannel=h5"
const lrsPackage = "com.netease.lrs.vivo"
const mlPackage = "com.c2vl.kgamebox"

type ret struct {
	Data struct {
		Games []appData `json:"games"`
	} `json:"data"`
}

type appData struct {
	PkgName string `json:"pkgName"`
	Apkurl  string `json:"apkurl"`
}

// GetRank .
func GetRank() (*types.RankInfo, error) {
	bs, err := util.HTTPGet(searchAPI, nil)
	if err != nil {
		return nil, err
	}
	ret := &ret{}
	err = json.Unmarshal(bs, ret)
	if err != nil {
		return nil, err
	}
	r := &types.RankInfo{}
	for i, v := range ret.Data.Games {
		if v.PkgName == lrsPackage {
			r.LrsRank = i + 1
		}
		if v.PkgName == mlPackage {
			r.MlRank = i + 1
			r.MlURL = v.Apkurl
		}
	}
	return r, nil
}
