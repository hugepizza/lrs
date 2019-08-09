package img

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
	fmt.Println(getToken())
}

func TestGetSimilar(t *testing.T) {
	fmt.Println(GetSimilar("http://p17.qhimg.com/t011cdd0e315df7d864.png"))
}
