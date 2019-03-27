package ipfs

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/cretz/tor-dht-poc/go-i2p-dht-poc/i2pdht/ipfs/websocket"
	gorillaws "github.com/gorilla/websocket"

	"github.com/whyrusleeping/mafmt"

	"github.com/eyedeekay/sam3"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"

	upgrader "github.com/libp2p/go-libp2p-transport-upgrader"
)

// impls libp2p's transport.Transport
type I2PTransport struct {
	samI2P   *sam3.SAM
	conf     *I2PTransportConf
	upgrader *upgrader.Upgrader

	dialerLock sync.Mutex
	i2pDialer  *sam3.StreamSession
	wsDialer   *gorillaws.Dialer
}

type I2PTransportConf struct {
	WebSocket bool
}

var EepMultiaddrFormat = mafmt.Base(ma.P_GARLIC64)
var I2PMultiaddrFormat = mafmt.Or(EepMultiaddrFormat, mafmt.TCP)

var _ transport.Transport = &I2PTransport{}

func NewI2PTransport(samI2P *sam3.SAM, conf *I2PTransportConf) func(*upgrader.Upgrader) *I2PTransport {
	return func(upgrader *upgrader.Upgrader) *I2PTransport {
		log.Printf("Creating transport with upgrader: %v", upgrader)
		if conf == nil {
			conf = &I2PTransportConf{}
		}
		return &I2PTransport{samI2P: samI2P, conf: conf, upgrader: upgrader}
	}
}

func (t *I2PTransport) Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (transport.Conn, error) {
	log.Printf("For peer ID %v, dialing %v", p, raddr)
	var addr string
	if garlicID, port, err := defaultAddrFormat.garlicInfo(raddr); err != nil {
		return nil, err
	} else {
		addr = fmt.Sprintf("%v.garlic:%v", garlicID, port)
	}
	// Init the dialers
	if err := t.initDialers(ctx); err != nil {
		log.Printf("Failed initializing dialers: %v", err)
		return nil, err
	}
	// Now dial
	var netConn net.Conn
	if t.wsDialer != nil {
		log.Printf("Dialing addr: ws://%v", addr)
		wsConn, _, err := t.wsDialer.Dial("ws://"+addr, nil)
		if err != nil {
			log.Printf("Failed dialing: %v", err)
			return nil, err
		}
		netConn = websocket.NewConn(wsConn, nil)
	} else {
		var err error
		if netConn, err = t.i2pDialer.DialContext(ctx, "tcp", addr); err != nil {
			log.Printf("Failed dialing: %v", err)
			return nil, err
		}
	}
	// Convert connection
	if manetConn, err := manet.WrapNetConn(netConn); err != nil {
		log.Printf("Failed wrapping the net connection: %v", err)
		return nil, err
	} else if conn, err := t.upgrader.UpgradeOutbound(ctx, t, manetConn, p); err != nil {
		log.Printf("Failed upgrading connection: %v", err)
		return nil, err
	} else {
		return conn, nil
	}
}

func (t *I2PTransport) initDialers(ctx context.Context) error {
	t.dialerLock.Lock()
	defer t.dialerLock.Unlock()
	// If already inited, good enough
	if t.i2pDialer != nil {
		return nil
	}
	//var err error
	keys, err := t.samI2P.NewKeys(sam3.Sig_EdDSA_SHA512_Ed25519)
	if err != nil {
		return err
	}
	if t.i2pDialer, err = t.samI2P.NewStreamSessionWithSignature("testdht", keys, []string{}, sam3.Sig_EdDSA_SHA512_Ed25519); err != nil {
		return fmt.Errorf("Failed creating samv3 StreamSession: %v", err)
	}
	// Create web socket dialer if needed
	if t.conf.WebSocket {
		t.wsDialer = &gorillaws.Dialer{
			NetDial:          t.i2pDialer.Dial,
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
		}
	}
	return nil
}

func (t *I2PTransport) CanDial(addr ma.Multiaddr) bool {
	log.Printf("Checking if can dial %v", addr)
	_, _, err := defaultAddrFormat.garlicInfo(addr)
	return err == nil
}

func (t *I2PTransport) Listen(laddr ma.Multiaddr) (transport.Listener, error) {
	// TODO: support a bunch of config options on this if we want
	log.Printf("Called listen for %v", laddr)
	if val, err := laddr.ValueForProtocol(ma.P_GARLIC64); err != nil {
		return nil, fmt.Errorf("Unable to get protocol value: %v", err)
	} else if val != "" {
		return nil, fmt.Errorf("Must be '/garlic64', got '/garlic64/%v'", val)
	}
	// Listen with version 3, wait 1 min for bootstrap
	//ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Minute)
	//defer cancelFn()
	garlic, err := t.i2pDialer.Listen() //ctx, &tor.ListenConf{Version3: true})
	if err != nil {
		log.Printf("Failed creating garlic service: %v", err)
		return nil, err
	}

	log.Printf("Listening on garlic: %v", t.i2pDialer.Addr().Base64())
	// Close it if there is another error in here
	defer func() {
		if err != nil {
			log.Printf("Failed listen after garlic creation: %v", err)
			garlic.Close()
		}
	}()

	// Return a listener
	manetListen := &manetListener{transport: t, garlic: garlic, listener: garlic}
	remoteaddr, _ := strconv.Atoi(garlic.To())
	addrStr := defaultAddrFormat.garlicAddr(garlic.Addr().String(), remoteaddr)
	if t.conf.WebSocket {
		addrStr += "/ws"
	}
	if manetListen.multiaddr, err = ma.NewMultiaddr(addrStr); err != nil {
		return nil, fmt.Errorf("Failed converting garlic address: %v", err)
	}
	// If it had websocket, we need to delegate to that
	if t.conf.WebSocket {
		if manetListen.listener, err = websocket.StartNewListener(garlic); err != nil {
			return nil, fmt.Errorf("Failed creating websocket: %v", err)
		}
	}

	log.Printf("Completed creating IPFS listener from garlic, addr: %v", manetListen.multiaddr)
	return manetListen.Upgrade(t.upgrader), nil
}

func (t *I2PTransport) Protocols() []int { return []int{ma.P_TCP, ma.P_GARLIC64} }
func (t *I2PTransport) Proxy() bool      { return true }

type manetListener struct {
	transport *I2PTransport
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
