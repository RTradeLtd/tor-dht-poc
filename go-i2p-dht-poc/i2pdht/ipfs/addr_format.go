package ipfs

import (
	"fmt"

	ma "github.com/multiformats/go-multiaddr"
)

type addrFormat interface {
	garlicInfo(addr ma.Multiaddr) (id string, port int, err error)
	garlicAddr(id string, port int) string
}

var defaultAddrFormat addrFormat = &addrFormatProtocol{}

type addrFormatProtocol struct{}

func (addrFormatProtocol) garlicInfo(addr ma.Multiaddr) (string, int, error) {
	if addr != nil {
		var garlicAddrStr string
		var err error
		if garlicAddrStr, err = addr.ValueForProtocol(ma.P_GARLIC32); err != nil {
			return "", -1, fmt.Errorf("Failed getting garlic info from %v: %v", addr, err)
		}
		return garlicAddrStr, 0, nil
	}
	return "", -1, nil //fmt.Errorf("Failed because garlic addr was nil")
}

func (addrFormatProtocol) garlicAddr(id string, unused int) string {
	return fmt.Sprintf("/garlic32/%v", id)
}
