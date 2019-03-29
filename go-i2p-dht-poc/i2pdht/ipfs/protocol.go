package ipfs

import (
	"strings"
)

func trim(b32 string) string {
	return strings.Replace(b32, ".b32.i2p", "", -1)
}
