package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/tealeg/xlsx"

	"github.com/robfig/cron"

	"github.com/sirupsen/logrus"

	"github.com/hugepizza/lrs/apps/baidu"
	"github.com/hugepizza/lrs/apps/huawei"
	"github.com/hugepizza/lrs/apps/meizu"
	"github.com/hugepizza/lrs/apps/oppo"
	"github.com/hugepizza/lrs/apps/threesixzero"
	"github.com/hugepizza/lrs/apps/tx"
	"github.com/hugepizza/lrs/apps/vivo"
	"github.com/hugepizza/lrs/apps/xiaomi"
	"github.com/hugepizza/lrs/types"
	"github.com/hugepizza/lrs/util"
	yaml "gopkg.in/yaml.v2"
)

var searchs = map[string]string{
	"baidu":  "https://m.baidu.com/s?ie=UTF-8&wd=%E7%8B%BC%E4%BA%BA%E6%9D%80",
	"sougou": "https://m.sogou.com/web/searchList.jsp?uID=sadasd%3D&v=5&dp=2&pid=sogou-waps-23asd&w=1283&t=1563724067454&s_t=1563724073043&s_from=result_up&htprequery=lrs&keyword=%E7%8B%BC%E4%BA%BA%E6%9D%80&pg=webSearchList&rcer=gNz_a8U1sUAKzX9o&s=%E6%90%9C%E7%B4%A2&suguuid=61sad537956-a974-45b6-af2c-97d33769d3b6&sugsuv=dasd&sugtime=1553724073043",
	"360":    "https://m.so.com/s?q=%E7%8B%BC%E4%BA%BA%E6%9D%80&src=suggest_history&sug_pos=0&sug=&srcg=home_next",
	"shenma": "https://m.sm.cn/s?q=%E7%8B%BC%E4%BA%BA%E6%9D%80&from=smor&safe=1&snum=1",

	"pc_baidu":  "https://www.baidu.com/s?ie=utf-8&f=3&rsv_bp=1&tn=baidu&wd=%E7%8B%BC%E4%BA%BA%E6%9D%80",
	"pc_sougou": "https://www.sogou.com/web?query=%E7%8B%BC%E4%BA%BA%E6%9D%80&_asf=www.sogou.com&_ast=&w=01019900&p=40040100&ie=utf8&from=index-nologin&s_from=index&sut=872&sst0=1563761139761&lkt=0%2C0%2C0&sugsuv=00563C730155C8CC5CEF84167D95B509&sugtime=1563761139761",
	"pc_360":    "https://www.so.com/s?ie=utf-8&fr=none&src=360sou_newhome&q=%E7%8B%BC%E4%BA%BA%E6%9D%80",
}

type config struct {
	Email struct {
		Sender    string   `yaml:"sender"`
		Code      string   `yaml:"code"`
		Server    string   `yaml:"server"`
		Port      int      `yaml:"port"`
		Receivers []string `yaml:"receivers"`
	} `yaml:"email"`
	BaiduAI struct {
		Key    string `yaml:"key"`
		Secret string `yaml:"secret"`
	} `yaml:"baiduai"`
}

var conf = &config{}

func init() {
	bs, err := ioutil.ReadFile("./conf.yaml")
	if err != nil {
		logrus.Fatal(err)
	}
	err = yaml.Unmarshal(bs, conf)
	if err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	var crond = cron.New()
	crond.AddFunc("0 30 18 * * ?", func() {
		run()
	})
	crond.AddFunc("0 30 8 * * ?", func() {
		run()
	})
	crond.AddFunc("0 0 23 * * ?", func() {
		run()
	})
	crond.Start()

	go run()

	select {}
}

func run() {
	files := searchShot(searchs)
	appRankFile := appRank()
	if appRankFile != "" {
		files = append(files, appRankFile)
	}
	defer func() {
		for _, f := range files {
			os.Remove(f)
		}
	}()
	shotAndsend(files)
	logrus.Info("send success")
}

func shotAndsend(attachFiles []string) error {
	return util.SendEmail(conf.Email.Sender,
		conf.Email.Code,
		conf.Email.Receivers,
		"狼人杀app搜索引擎关键字监控",
		attachFiles,
		conf.Email.Server,
		conf.Email.Port)
}

func searchShot(data map[string]string) []string {
	files := []string{}
	for k, v := range data {
		fn := filepath.Join(os.TempDir(), fmt.Sprintf("%s.png", k))
		err := util.Shot(fn, v)
		if err != nil {
			continue
		}
		files = append(files, fn)
	}
	return files
}

func appRank() string {
	ranks := make(map[string]*types.RankInfo)
	logrus.Info("request 华为")
	rank, err := huawei.GetRank()
	if err != nil {
		logrus.Error("华为", err)
	} else {
		ranks["华为"] = rank
	}
	logrus.Info("request 小米")
	rank2, err := xiaomi.GetRank()
	if err != nil {
		logrus.Error("小米", err)
	} else {
		ranks["小米"] = rank2
	}
	logrus.Info("request 百度")
	rank3, err := baidu.GetRank()
	if err != nil {
		logrus.Error("百度", err)
	} else {
		ranks["百度"] = rank3
	}
	logrus.Info("request 应用宝")
	rank4, err := tx.GetRank()
	if err != nil {
		logrus.Error("应用宝", err)
	} else {
		ranks["应用宝"] = rank4
	}
	logrus.Info("request oppo")
	rank5, err := oppo.GetRank()
	if err != nil {
		logrus.Error("oppo", err)
	} else {
		ranks["oppo"] = rank5
	}
	logrus.Info("request 360")
	rank6, err := threesixzero.GetRank()
	if err != nil {
		logrus.Error("360", err)
	} else {
		ranks["360"] = rank6
	}
	logrus.Info("request 魅族")
	rank7, err := meizu.GetRank()
	if err != nil {
		logrus.Error("魅族", err)
	} else {
		ranks["魅族"] = rank7
	}
	logrus.Info("request vivo")
	rank8, err := vivo.GetRank()
	if err != nil {
		logrus.Error("vivo", err)
	} else {
		ranks["vivo"] = rank8
	}
	return genAppRankXlsx(ranks)
}

func genAppRankXlsx(data map[string]*types.RankInfo) string {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row := sheet.AddRow()
	cell0 := row.AddCell()
	cell0.Value = "平台"
	cell1 := row.AddCell()
	cell1.Value = "官狼排名"
	cell2 := row.AddCell()
	cell2.Value = "红狼排名"
	cell3 := row.AddCell()
	cell3.Value = "红狼下载"
	for k, v := range data {
		row := sheet.AddRow()
		cell0 := row.AddCell()
		cell0.Value = k
		cell1 := row.AddCell()
		if v.LrsRank > 0 {
			cell1.Value = fmt.Sprintf("%d", v.LrsRank)
		} else {
			if k == "魅族" {
				cell1.Value = "无法追踪"
			} else {
				cell1.Value = "无"
			}
		}
		cell2 := row.AddCell()
		if v.MlRank > 0 {
			cell2.Value = fmt.Sprintf("%d", v.MlRank)
		} else {
			cell2.Value = "无"
		}
		cell3 := row.AddCell()
		cell3.Value = v.MlURL
	}
	path := filepath.Join(os.TempDir(), time.Now().Format("lrs_2006010215")+".xlsx")
	err = file.Save(path)
	if err != nil {
		return ""
	}
	return path
}
