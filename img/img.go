package img

import (
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"os"

	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"strings"

	"github.com/Comdex/imgo"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/draw"
	"golang.org/x/image/riff"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

func stdHash(w, h int) string {
	tf, err := ioutil.TempFile("", "*"+".png")
	if err != nil {
		logrus.Error(err)
		return ""
	}
	// defer os.Remove(tf.Name())
	fmt.Println(tf.Name())
	err = scale("std.png", tf.Name(), w, h, true, false)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	std, err := imgo.GetFingerprint(tf.Name())
	if err != nil {
		logrus.Error(err)
		return ""
	}
	return std
}

func getHash(src string) (string, int, int) {
	tf, err := ioutil.TempFile("", "")
	if err != nil {
		logrus.Error(err)
		return "", 0, 0
	}
	defer os.Remove(tf.Name())
	resp, err := http.Get(src)
	if err != nil {
		logrus.Error(err)
		return "", 0, 0
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return "", 0, 0
	}
	_, err = tf.Write(bs)
	if err != nil {
		logrus.Error(err)
		return "", 0, 0
	}
	w, h, err := imageSizeWithFormat(tf.Name())
	if err != nil {
		logrus.Error(err)
		return "", 0, 0
	}
	hash, err := imgo.GetFingerprint(tf.Name())
	if err != nil {
		logrus.Error(err)
		return "", 0, 0
	}
	return hash, w, h
}

// GetDiff .
func GetDiff(src string) int {
	input, w, h := getHash(src)
	std := stdHash(w, h)
	fmt.Println(input)
	fmt.Println(std)
	if len(std) != 64 || len(input) != 64 {
		return 64
	}
	diff := 0
	for index := 0; index < 64; index++ {
		if std[index] != input[index] {
			diff++
		}
	}
	return diff
}

func scale(src, dst string, w, h int, equalRate bool, cut bool) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}

	ext := filepath.Ext(src)
	var imgSrc image.Image

	switch ext {
	case ".jpg", ".jpeg":
		imgSrc, err = jpeg.Decode(f)
	case ".png":
		imgSrc, err = png.Decode(f)
	case ".gif":
		imgSrc, err = gif.Decode(f)
	default:
		err = fmt.Errorf("unknown input image type")
	}

	f.Close()

	if err != nil {
		return err
	}
	if dst == "" || dst == src {
		dst = src
	}
	f, err = os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	ext = filepath.Ext(dst)

	switch ext {
	case ".jpg", ".jpeg":
		err = nil
	case ".png":
		err = nil
	case ".gif":
		err = nil
	default:
		err = fmt.Errorf("unknown output image type")
	}
	if err != nil {
		return err
	}

	// 计算缩放比例，当启用裁剪时按小比率，否则为大比率
	if cut {
		// 当启用裁剪时，一定使用等比
		equalRate = true
	}
	var rate float64
	srcBound := imgSrc.Bounds()
	rateX := float64(srcBound.Dx()) / float64(w)
	rateY := float64(srcBound.Dy()) / float64(h)
	if rateX == rateY {
		rate = rateX // 正好等比
	} else {
		if rateX > rateY {
			if cut {
				rate = rateY
			} else {
				rate = rateX
			}
		} else {
			if cut {
				rate = rateX
			} else {
				rate = rateY
			}
		}
	}

	// 计算目标尺寸，启用等比时，目标尺寸不等于设定尺寸，所以要先计算比率
	dstW := int(float64(srcBound.Dx()) / rate)
	dstH := int(float64(srcBound.Dy()) / rate)
	// 如果没启用等比，一定没启用裁剪，直接赋值，不考虑比率问题
	if equalRate == false {
		dstW = w
		dstH = h
	}

	// 首次处理，先缩放
	imgDst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	draw.CatmullRom.Scale(imgDst, imgDst.Bounds(), imgSrc, imgSrc.Bounds(), draw.Src, nil)

	// 不等比，或者正好比率相同也无需裁剪，直接保存
	if equalRate && rateX != rateY {
		// 裁剪图片，输出目标尺寸，截取中间部分
		imgDst = imgDst.SubImage(image.Rectangle{
			Min: image.Point{(dstW - w) / 2, (dstH - h) / 2},
			Max: image.Point{w + (dstW-w)/2, h + (dstH-h)/2},
		}).(*image.RGBA)
	}

	// 导出
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(f, imgDst, &jpeg.Options{Quality: 95})
	case ".png":
		encoder := png.Encoder{
			CompressionLevel: png.BestCompression,
		}
		err = encoder.Encode(f, imgDst)
	case ".gif":
		err = gif.Encode(f, imgDst, nil)
	default:
		err = fmt.Errorf("unknown output image type")
	}
	return err
}

func imageSizeWithFormat(file string) (width int, height int, err error) {
	f, err := os.Open(file)
	if err != nil {
		return 0, 0, err
	}
	format, err := imageFormatDiscernment(f)
	if err != nil {
		return 0, 0, err
	}
	var img image.Image
	switch format {
	case "JPEG":
		img, err = jpeg.Decode(f)
	case "PNG":
		img, err = png.Decode(f)
	case "GIF":
		img, err = gif.Decode(f)
	case "TIFF":
		img, err = tiff.Decode(f)
	case "WEBP":
		img, err = webp.Decode(f)
	default:
		err = errors.New("not support")
	}
	f.Close()
	if err != nil {
		return 0, 0, err
	}

	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), nil
}

func imageFormatDiscernment(f *os.File) (string, error) {
	head := make([]byte, 8)
	_, err := f.ReadAt(head, 0)
	if err != nil {
		return "", err
	}
	headStr := strings.ToUpper(hex.EncodeToString(head))
	_, err = f.Seek(0, 0)
	if err != nil {
		return "", err
	}

	format := ""
	if strings.HasPrefix(headStr, "FFD8") {
		format = "JPEG"
	} else if strings.HasPrefix(headStr, "89504E470D0A1A0A") {
		format = "PNG"
	} else if strings.HasPrefix(headStr, "474946383961") || strings.HasPrefix(headStr, "474946383761") {
		format = "GIF"
	} else if strings.HasPrefix(headStr, "4D4D") || strings.HasPrefix(headStr, "4949") {
		format = "TIFF"
	} else if strings.HasPrefix(headStr, "52494646") {
		fccWEBP := riff.FourCC{'W', 'E', 'B', 'P'}
		formType, _, err := riff.NewReader(f)
		f.Seek(0, 0)
		if err == nil {
			if formType == fccWEBP {
				format = "WEBP"
			}
		} else {
			err = errors.New("not support")
		}
	} else {
		err = errors.New("not support")
	}

	if err != nil {
		return "", err
	}
	return format, nil
}
