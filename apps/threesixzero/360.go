package threesixzero

import (
	"bytes"
	"strings"

	"github.com/hugepizza/lrs/util"

	"github.com/PuerkitoBio/goquery"
	"github.com/hugepizza/lrs/types"
)

const searchAPI = "http://zhushou.360.cn/search/index/?kw=%E7%8B%BC%E4%BA%BA%E6%9D%80"
const lrsAPPID = "3906146"
const mlAPPID = "3583525"

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
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	r := &types.RankInfo{}
	doc.Find(".SeaCon").Find("ul").Find("li").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("h3").Find("a").Attr("href")
		if strings.Contains(href, lrsAPPID) {
			r.LrsRank = i + 1
		}
		if strings.Contains(href, mlAPPID) {
			r.MlRank = i + 1
		}
	})
	return r, nil
}
