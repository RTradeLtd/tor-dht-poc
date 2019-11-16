package ipfs

import (
	"fmt"
	"github.com/RTradeLtd/go-garlic-tcp-transport"
	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"
	"math/rand"
)

type I2PTransportConf struct {
	i2pkeys.I2PKeys
	WebSocket bool
}

func NewID() string {
	b := make([]byte, 12)
	for i := range b {
		b[i] = "abcdefghijklmnopqrstuvwxyz"[rand.Intn(len("abcdefghijklmnopqrstuvwxyz"))]
	}
	return "dht" + string(b)
}

func NewI2PTransport() (*i2ptcp.GarlicTCPTransport, error) {
	// Create the host with only the i2p transport
	samI2P, err := sam3.NewSAM("127.0.0.1:7656")
	if err != nil {

		return nil, err
	}
	defer samI2P.Close()
	//t.debugf("Creating host")
	//k, err := samI2P.NewKeys(sam3.Sig_EdDSA_SHA512_Ed25519)
	//if err != nil {
	//return nil, fmt.Errorf("Failed generating I2P keys: %v",err)
	//}
	garlicTransport, err := i2ptcp.NewGarlicTCPTransportFromOptions()
	if err != nil {
		return nil, fmt.Errorf("Failed setting up garlic transport: %v", err)
	}
	return garlicTransport, nil
}
