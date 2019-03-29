package ipfs

import (
	"log"
	"strings"
)

func trim(b32 string) string {
	log.Println("BASE32 WAS", strings.Replace(b32, ".b32.i2p", "", -1))
	return strings.Replace(b32, ".b32.i2p", "", -1)
}
