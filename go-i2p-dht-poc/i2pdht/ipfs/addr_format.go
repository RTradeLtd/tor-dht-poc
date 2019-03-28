package ipfs

import (
	"fmt"

	ma "github.com/multiformats/go-multiaddr"
)

type addrFormat interface {
	garlicInfo(addr ma.Multiaddr) (id string, port int, err error)
	garlicAddr(id string, port int) string
}

var defaultAddrFormat addrFormat = addrFormatProtocol{}

type addrFormatProtocol struct{}

func (addrFormatProtocol) garlicInfo(addr ma.Multiaddr) (string, int, error) {
	var garlicAddrStr string
	var err error
	if garlicAddrStr, err = addr.ValueForProtocol(ma.P_GARLIC64); err != nil {
		return "", -1, fmt.Errorf("Failed getting garlic info from %v: %v", addr, err)
	}
	return garlicAddrStr, 0, nil

}

func (addrFormatProtocol) garlicAddr(id string, unused int) string {
	return fmt.Sprintf("/garlic64/%v", id)
}
