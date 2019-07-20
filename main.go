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

var searchs = map[string]string{
	"baidu":  "https://m.baidu.com/s?ie=UTF-8&wd=%E7%8B%BC%E4%BA%BA%E6%9D%80",
	"sougou": "https://wap.sogou.com/web/searchList.jsp?keyword=狼人杀",
	"360":    "https://m.so.com/s?q=%E7%8B%BC%E4%BA%BA%E6%9D%80&src=suggest_history&sug_pos=0&sug=&srcg=home_next",
	"shenma": "https://m.sm.cn/s?q=%E7%8B%BC%E4%BA%BA%E6%9D%80&from=smor&safe=1&snum=1",
}

func main() {
	var crond = cron.New()
	crond.AddFunc("0 0 9 * * *", func() {
		sendLrs()
	})
	go sendLrs()
	select {}
}

func sendLrs() {
	if err := shot(); err != nil {
		log.Println(err)
	}
	if err := sendEmail(); err != nil {
		log.Println(err)
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
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(map[string]interface{}{
			"user-agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1",
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
