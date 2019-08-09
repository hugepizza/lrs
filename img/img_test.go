package img

import (
	"fmt"
	"testing"
)

func TestGetHash(t *testing.T) {
	fmt.Println(getHash("http://p17.qhimg.com/t011cdd0e315df7d864.png"))
}

func TestStdHash(t *testing.T) {
	fmt.Println(stdHash(72, 72))
	fmt.Println(getHash("http://p17.qhimg.com/t011cdd0e315df7d864.png"))
}

func TestGetDiff(t *testing.T) {
	GetDiff("http://p17.qhimg.com/t011cdd0e315df7d864.png")
}
