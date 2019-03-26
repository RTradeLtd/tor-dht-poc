package ipfs

import (
	//"encoding/base32"
	"fmt"

	ma "github.com/multiformats/go-multiaddr"
)

var garlicListenAddr ma.Multiaddr

func init() {
    var err error
	// Add the listen protocol
	if garlicListenAddr, err = ma.NewMultiaddr("/garlic64"); err != nil {
		panic(fmt.Errorf("Failed creating garlic64 addr: %v", err))
	}
}
