package i2pdht

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

import (
	"github.com/cretz/bine/torutil"
	"github.com/eyedeekay/sam3"
)

type Impl interface {
	ApplyDebugLogging()
	RawStringDataID(id []byte) (string, error)
	NewDHT(ctx context.Context, conf *DHTConf) (DHT, error)
}

type DHTConf struct {
	I2P            *sam3.SAM
	BootstrapPeers []*PeerInfo
	ClientOnly     bool
	Verbose        bool
}

type DHT interface {
	io.Closer

	PeerInfo() *PeerInfo
	Provide(ctx context.Context, id []byte) error
	FindProviders(ctx context.Context, id []byte, maxCount int) ([]*PeerInfo, error)
}

type PeerInfo struct {
	ID string
	// May be empty string if not listening
	EepServiceID string
	// Invalid value if EepServiceID is empty
	EepPort int
}

func (p *PeerInfo) String() string {
	return fmt.Sprintf("/garlic32/%v/tcp/%v/%v", p.EepServiceID, p.EepPort, p.ID)
}

func NewPeerInfo(str string) (*PeerInfo, error) {
	if garlic, id, ok := torutil.PartitionString(str, '/'); !ok {
		return nil, fmt.Errorf("Missing ID portion")
	} else if garlicID, portStr, ok := torutil.PartitionString(garlic, ':'); !ok {
		return nil, fmt.Errorf("Missing garlic port")
	} else {
		ret := &PeerInfo{ID: id, EepServiceID: garlicID}
		ret.EepPort, _ = strconv.Atoi(portStr)
		return ret, nil
	}
}
