package oppo

import (
	"encoding/json"

	"github.com/hugepizza/lrs/types"
	"github.com/hugepizza/lrs/util"
)

const searchAPI = "https://istore.oppomobile.com/search/v1/search?tabId=&start=0&size=10&keyword=%E7%8B%BC%E4%BA%BA%E6%9D%80"
const lrsPackageName = "com.netease.lrs"
const mlPackageName = "com.c2vl.kgamebox"

type ret struct {
	Cards []card `json:"cards"`
}

type card struct {
	App  *app   `json:"app"`
	Apps []*app `json:"apps"`
}

type app struct {
	PkgName string `json:"pkgName"`
	URL     string `json:"url"`
}

// GetRank .
func GetRank() (*types.RankInfo, error) {
	header := map[string][]string{
		"Accept-Encoding": []string{"gzip"},
		"User-Agent":      []string{"nubia%2FNX523J_V1%2F19%2F4.4.2%2F0%2F2%2F2101%2F7103"},
		"locale":          []string{"zh-CN;CN"},
		"appversion":      []string{"7.1.0"},
		"appid":           []string{"nubia#001#CN"},
		"pid":             []string{"001"},
		"sg":              []string{"dfad3bc14acd3d7e15371f1927e928c3ae4c9e72"},
		"ct":              []string{"1565074126589"},
		"nw":              []string{"1"},
		"ocs":             []string{"nubia%2FNX523J_V1%2F19%2F4.4.2%2F0%2F2%2Fsamsung-user+4.4.2+KOT49H+3.8.017.0602+release-keys%2F7103"},
		"country":         []string{"cn"},
		"sign":            []string{"f65a1596282c77b069e67e817b34c61f"},
		"id":              []string{"864394010201319"},
		"t":               []string{"1565074126579"},
		"token":           []string{"1"},
		"oak":             []string{"cdb09c43063ea6bb"},
		"gid":             []string{"abf7f743-2437-44b5-86c3-3d1ccd683972"},
		"rsq":             []string{"19843"},
		"ch":              []string{"2101"},
		"Accept":          []string{"application/json; charset=UTF-8"},
		"Host":            []string{"istore.oppomobile.com"},
		"Connection":      []string{"Keep-Alive"},
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
	apps := []*app{}
	for _, v := range ret.Cards {
		if v.App != nil {
			apps = append(apps, v.App)
		}
		if v.Apps != nil {
			if err == nil {
				apps = append(apps, v.Apps...)
			}
		}
	}
	for i, v := range apps {
		if v.PkgName == lrsPackageName {
			r.LrsRank = i + 1
		}
		if v.PkgName == mlPackageName {
			r.MlRank = i + 1
			r.MlURL = v.URL
		}
	}
	return r, nil
}
