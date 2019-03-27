package ipfs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cretz/bine/torutil"
	ma "github.com/multiformats/go-multiaddr"
	madns "github.com/multiformats/go-multiaddr-dns"
)

type addrFormat interface {
	garlicInfo(addr ma.Multiaddr) (id string, port int, err error)
	garlicAddr(id string, port int) string
}

// var defaultAddrFormat addrFormat = addrFormatProtocol{}
var defaultAddrFormat addrFormat = addrFormatDns{}

// In the form /garlic/<garlic-id>:<port>
type addrFormatProtocol struct{}

func (addrFormatProtocol) garlicInfo(addr ma.Multiaddr) (string, int, error) {
	if garlicAddrStr, err := addr.ValueForProtocol(ma.P_ONION); err != nil {
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
	return fmt.Sprintf("/garlic/%v:%v", id, port)
}

// In the form /dns4/<garlic-id>.garlic/tcp/<port>
type addrFormatDns struct{}

func (addrFormatDns) garlicInfo(addr ma.Multiaddr) (string, int, error) {
	if addrPieces := ma.Split(addr); len(addrPieces) < 2 {
		return "", -1, fmt.Errorf("Invalid pieces: %v", addrPieces)
	} else if garlicAddrStr, err := addrPieces[0].ValueForProtocol(madns.Dns4Protocol.Code); err != nil {
		return "", -1, fmt.Errorf("Can't get garlic part of %v: %v", addr, err)
	} else if !strings.HasSuffix(garlicAddrStr, ".garlic") {
		return "", -1, fmt.Errorf("Invalid garlic addr: %v", garlicAddrStr)
	} else if portStr, err := addrPieces[1].ValueForProtocol(ma.P_TCP); err != nil {
		return "", -1, fmt.Errorf("Can't get port part of %v: %v", addr, err)
	} else if port, portErr := strconv.Atoi(portStr); portErr != nil {
		return "", -1, fmt.Errorf("Invalid port '%v': %v", portStr, portErr)
	} else {
		return garlicAddrStr[:len(garlicAddrStr)-6], port, nil
	}
}

func (addrFormatDns) garlicAddr(id string, port int) string {
	return fmt.Sprintf("/dns4/%v.garlic/tcp/%v", id, port)
}
