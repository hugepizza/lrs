package util

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}
var cli = &http.Client{
	Timeout:   time.Second * 10,
	Transport: tr,
}

// HTTPGet .
func HTTPGet(url string, header map[string][]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = header
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
