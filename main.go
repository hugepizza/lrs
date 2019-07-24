package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/robfig/cron"
	"gopkg.in/gomail.v1"
)

const mAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1"

const pcAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36"

var searchs = map[string]string{
	"baidu":  "https://m.baidu.com/s?ie=UTF-8&wd=%E7%8B%BC%E4%BA%BA%E6%9D%80",
	"sougou": "https://m.sogou.com/web/searchList.jsp?uID=sadasd%3D&v=5&dp=2&pid=sogou-waps-23asd&w=1283&t=1563724067454&s_t=1563724073043&s_from=result_up&htprequery=lrs&keyword=%E7%8B%BC%E4%BA%BA%E6%9D%80&pg=webSearchList&rcer=gNz_a8U1sUAKzX9o&s=%E6%90%9C%E7%B4%A2&suguuid=61sad537956-a974-45b6-af2c-97d33769d3b6&sugsuv=dasd&sugtime=1553724073043",
	"360":    "https://m.so.com/s?q=%E7%8B%BC%E4%BA%BA%E6%9D%80&src=suggest_history&sug_pos=0&sug=&srcg=home_next",
	"shenma": "https://m.sm.cn/s?q=%E7%8B%BC%E4%BA%BA%E6%9D%80&from=smor&safe=1&snum=1",

	"pc_baidu":  "https://www.baidu.com/s?ie=utf-8&f=3&rsv_bp=1&tn=baidu&wd=%E7%8B%BC%E4%BA%BA%E6%9D%80",
	"pc_sougou": "https://www.sogou.com/web?query=%E7%8B%BC%E4%BA%BA%E6%9D%80&_asf=www.sogou.com&_ast=&w=01019900&p=40040100&ie=utf8&from=index-nologin&s_from=index&sut=872&sst0=1563761139761&lkt=0%2C0%2C0&sugsuv=00563C730155C8CC5CEF84167D95B509&sugtime=1563761139761",
	"pc_360":    "https://www.so.com/s?ie=utf-8&fr=none&src=360sou_newhome&q=%E7%8B%BC%E4%BA%BA%E6%9D%80",
}

func main() {
	var crond = cron.New()
	crond.AddFunc("0 30 6 * * ?", func() {
		sendLrs()
	})
	crond.AddFunc("0 30 8 * * ?", func() {
		sendLrs()
	})
	crond.AddFunc("0 0 23 * * ?", func() {
		sendLrs()
	})
	crond.Start()
	go sendLrs()
	select {}
}

func sendLrs() {
	log.Printf("sending to %s \n", os.Getenv("LRS_SEND_LIST"))
	if err := shot(); err != nil {
		log.Println(err)
		return
	}
	if err := sendEmail(); err != nil {
		log.Println(err)
		return
	}
	log.Printf("send screenshot success at %s \n", time.Now().Format("2006/01/02 15:04:05"))
}

func shot() error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	for name, url := range searchs {
		var buf []byte
		if err := chromedp.Run(ctx, fullScreenshot(url, 100, &buf)); err != nil {
			return err
		}
		if err := ioutil.WriteFile(name+".png", buf, 0644); err != nil {
			return err
		}
	}
	return nil
}

func sendEmail() error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", os.Getenv("MY_QQMAIL_ADDR"))
	msg.SetHeader("To", strings.Split(os.Getenv("LRS_SEND_LIST"), "|")...)
	msg.SetHeader("Subject", fmt.Sprintf("%s 狼人杀搜索引擎关键字监控\n", time.Now().Format("2006年01月02日15点04分")))

	for k := range searchs {
		f, err := gomail.OpenFile(k + ".png")
		if err != nil {
			return err
		}
		msg.Attach(f)
	}

	mailer := gomail.NewMailer("smtp.qq.com", os.Getenv("MY_QQMAIL_ADDR"), os.Getenv("MY_QQMAIL_CODE"), 465)
	if err := mailer.Send(msg); err != nil {
		return err
	}
	return nil
}

func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	userAgent := mAgent
	if strings.HasPrefix(urlstr, "https://www") {
		userAgent = pcAgent
	}
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(map[string]interface{}{
			"user-agent": userAgent,
		})),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
