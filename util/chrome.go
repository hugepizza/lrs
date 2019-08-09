package util

import (
	"context"
	"io/ioutil"
	"math"
	"strings"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const mAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1"

const pcAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36"

// Shot 截图
func Shot(fileName, url string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var buf []byte
	if err := chromedp.Run(ctx, fullScreenshot(url, 100, &buf)); err != nil {
		return err
	}
	if err := ioutil.WriteFile(fileName, buf, 0644); err != nil {
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
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

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
