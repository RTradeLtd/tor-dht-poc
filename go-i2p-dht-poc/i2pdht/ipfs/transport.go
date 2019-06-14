package ipfs

import (
	/*	"context"
		"fmt"
		"log"
	*/
	"math/rand"
	"net"
	/*
		"net/http"
		"sync"
		"time"

		"github.com/cretz/tor-dht-poc/go-i2p-dht-poc/i2pdht/ipfs/websocket"
		gorillaws "github.com/gorilla/websocket"
	*/
	"github.com/whyrusleeping/mafmt"

	"github.com/RTradeLtd/go-garlic-tcp-transport"
	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"
	//	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"

	upgrader "github.com/libp2p/go-libp2p-transport-upgrader"
)

type I2PTransportConf struct {
	WebSocket bool
	keys      i2pkeys.I2PKeys
}

var EepMultiaddrFormat = mafmt.Base(ma.P_GARLIC32)
var I2PMultiaddrFormat = mafmt.Or(EepMultiaddrFormat, mafmt.TCP)

var _ transport.Transport = &i2ptcp.GarlicTCPTransport{} //&I2PTransport{}

func NewID() string {
	b := make([]byte, 12)
	for i := range b {
		b[i] = "abcdefghijklmnopqrstuvwxyz"[rand.Intn(len("abcdefghijklmnopqrstuvwxyz"))]
	}
	return "dht" + string(b)
}

type manetListener struct {
	transport *i2ptcp.GarlicTCPTransport
	garlic    *sam3.StreamListener
	multiaddr ma.Multiaddr
	listener  net.Listener
}

func (m *manetListener) Accept() (manet.Conn, error) {
	if c, err := m.listener.Accept(); err != nil {
		return nil, err
	} else {
		ret := &manetConn{Conn: c, localMultiaddr: m.multiaddr}
		if ret.remoteMultiaddr, err = manet.FromNetAddr(c.RemoteAddr()); err != nil {
			return nil, err
		}
		return ret, nil
	}
}

func (m *manetListener) Close() error            { return m.garlic.Close() }
func (m *manetListener) Addr() net.Addr          { return m.garlic.Addr() }
func (m *manetListener) Multiaddr() ma.Multiaddr { return m.multiaddr }
func (m *manetListener) Upgrade(u *upgrader.Upgrader) transport.Listener {
	return u.UpgradeListener(m.transport, m)
}

type manetConn struct {
	net.Conn
	localMultiaddr  ma.Multiaddr
	remoteMultiaddr ma.Multiaddr
}

func (m *manetConn) LocalMultiaddr() ma.Multiaddr  { return m.localMultiaddr }
func (m *manetConn) RemoteMultiaddr() ma.Multiaddr { return m.remoteMultiaddr }
