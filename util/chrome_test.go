package util

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestShot(t *testing.T) {
	urls := []string{
		"https://www.idongde.com/1f152bC0eE0e1a6F.shtml",
		"https://m.idongde.com/1f152bC0eE0e1a6F.shtml",
	}
	for i, url := range urls {
		err := Shot(filepath.Join(os.TempDir(), fmt.Sprintf("%d.png", i)), url)
		if err != nil {
			t.Fail()
		}
	}
}
