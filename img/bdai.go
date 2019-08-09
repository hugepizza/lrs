package img

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	tokenAPI  = "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s"
	searchAPI = "https://aip.baidubce.com/rest/2.0/realtime_search/same_hq/search"
)

type token struct {
	AccessToken string `json:"access_token"`
}
type searchResult struct {
	Result []struct {
		Score float64 `json:"score"`
	} `json:"result"`
}

func getToken() string {
	resp, err := http.Get(fmt.Sprintf(tokenAPI, "G89CK4qdMjOOdvzFQkG7zS85", "XUkF3oEl213hXYEtivhZOts7CgjMZTxy"))
	if err != nil {
		logrus.Error(err)
		return ""
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	token := &token{}
	err = json.Unmarshal(bs, token)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	return token.AccessToken
}

// GetSimilar 获取近似分
func GetSimilar(src string) float64 {
	data := url.Values{}
	data.Add("url", src)
	resp, err := http.Post(searchAPI+"?access_token="+getToken(), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		logrus.Error(err)
		return 0
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return 0
	}
	fmt.Println(string(bs))
	res := &searchResult{}
	err = json.Unmarshal(bs, res)
	if err != nil {
		logrus.Error(err)
		return 0
	}
	if len(res.Result) > 0 {
		return res.Result[0].Score
	}
	return 0
}
