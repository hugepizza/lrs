package xiaomi

import (
	"encoding/json"
	"fmt"

	"github.com/hugepizza/lrs/img"
	"github.com/hugepizza/lrs/util"

	"github.com/hugepizza/lrs/types"
)

const searchAPI = "http://m.app.mi.com/searchapi?keywords=%E7%8B%BC%E4%BA%BA%E6%9D%80&pageIndex=0&pageSize=20"
const lrsAPPID = 509823
const mlAPPID = 0

type ret struct {
	Data []appData `json:"data"`
}

type appData struct {
	AppID int    `json:"appID"`
	Icon  string `json:"icon"`
}

// RankInfo .
type RankInfo struct {
	LrsRank int
	MlRank  int
	MlURL   string
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
	for i, v := range ret.Data {
		if v.AppID == lrsAPPID {
			r.LrsRank = i + 1
		}
		same := img.GetSimilar(v.Icon)
		fmt.Println(same)
		if same > 0.8 {
			r.MlRank = i + 1
		}
		if r.LrsRank > 0 && r.MlRank > 0 {
			break
		}
	}
	return r, nil
}
