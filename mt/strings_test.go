package mt

import (
	"encoding/hex"
	"fmt"
	"log"
	"testing"
)

func TestToInt(t *testing.T) {
	log.Println(ToInt(""))
}

func TestEUCKR(t *testing.T) {
	fmt.Println(hex.EncodeToString([]byte("\xbe\xc6\xb8\xa7\xb4\xd9\xbf\xee \xbf\xec\xb8\xae\xb8\xbb")))
	fmt.Println(hex.EncodeToString(UTF8toEUCKR("아름다운 우리말")))

	fmt.Println(hex.EncodeToString(UTF8toEUCKR("아름다운 हि리말")))
}
