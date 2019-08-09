package util

import (
	"os"
	"testing"
)

func TestSendEmail(t *testing.T) {
	file := []string{
		"/Users/wll/go/src/github.com/hugepizza/lrs/go.mod",
		"/Users/wll/go/src/github.com/hugepizza/lrs/go.sum",
	}
	err := SendEmail("984373330@qq.com", os.Getenv("MY_QQMAIL_CODE"), []string{"984373330@qq.com"}, "test", file, "smtp.qq.com", 465)
	if err != nil {
		t.Fail()
	}
}
