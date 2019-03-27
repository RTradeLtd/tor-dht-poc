package ipfs

import (
	"fmt"
	"strconv"

	"github.com/cretz/bine/torutil"
	ma "github.com/multiformats/go-multiaddr"
)

type addrFormat interface {
	garlicInfo(addr ma.Multiaddr) (id string, port int, err error)
	garlicAddr(id string, port int) string
}

var defaultAddrFormat addrFormat = addrFormatProtocol{}

type addrFormatProtocol struct{}

func (addrFormatProtocol) garlicInfo(addr ma.Multiaddr) (string, int, error) {
	if garlicAddrStr, err := addr.ValueForProtocol(ma.P_GARLIC64); err != nil {
		return "", -1, fmt.Errorf("Failed getting garlic info from %v: %v", addr, err)
	} else if id, portStr, ok := torutil.PartitionString(garlicAddrStr, ':'); !ok {
		return "", -1, fmt.Errorf("Missing port on %v", garlicAddrStr)
	} else if port, portErr := strconv.Atoi(portStr); portErr != nil {
		return "", -1, fmt.Errorf("Invalid port '%v': %v", portStr, portErr)
	} else {
		return id, port, nil
	}
}

func (addrFormatProtocol) garlicAddr(id string, port int) string {
	return fmt.Sprintf("/garlic64/%v/tcp/%v", id, port)
}
