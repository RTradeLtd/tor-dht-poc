package ipfs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cretz/tor-dht-poc/go-i2p-dht-poc/i2pdht"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	//dht "github.com/RTradeLtd/go-libp2p-i2p-netdb"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"

	ma "github.com/multiformats/go-multiaddr"

	"github.com/eyedeekay/sam3"
	addr "github.com/ipfs/go-ipfs-addr"
)

type i2pDHT struct {
	debug    bool
	i2p      *sam3.SAM
	ipfsHost host.Host
	ipfsDHT  *dht.IpfsDHT
	peerInfo *i2pdht.PeerInfo
}

func (t *i2pDHT) Close() (err error) {
	if t.ipfsDHT != nil {
		err = t.ipfsDHT.Close()
	}
	if t.ipfsHost != nil {
		if hostCloseErr := t.ipfsHost.Close(); hostCloseErr != nil {
			// Just overwrite
			err = hostCloseErr
		}
	}
	return
}

func (t *i2pDHT) PeerInfo() *i2pdht.PeerInfo { return t.peerInfo }

func (t *i2pDHT) Provide(ctx context.Context, id []byte) error {
	if cid, err := ipfsImpl.hashedCID(id); err != nil {
		return err
	} else {
		t.debugf("Providing CID: %v", cid)
		return t.ipfsDHT.Provide(ctx, *cid, true)
	}
}

func (t *i2pDHT) FindProviders(ctx context.Context, id []byte, maxCount int) ([]*i2pdht.PeerInfo, error) {
	cid, err := ipfsImpl.hashedCID(id)
	if err != nil {
		return nil, err
	}
	t.debugf("Finding providers for CID: %v", cid)
	ret := []*i2pdht.PeerInfo{}
	for p := range t.ipfsDHT.FindProvidersAsync(ctx, *cid, maxCount) {
		if info, err := t.makePeerInfo(p.ID, p.Addrs[0]); err != nil {
			// TODO: warn instead?
			return nil, fmt.Errorf("Failed parsing '%v': %v", p, err)
		} else {
			ret = append(ret, info)
		}
	}
	return ret, ctx.Err()
}

func (t *i2pDHT) debugf(format string, args ...interface{}) {
	if t.debug {
		log.Printf("[DEBUG] "+format, args...)
	}
}

func (t *i2pDHT) applyPeerInfo() error {
	if listenAddrs := t.ipfsHost.Network().ListenAddresses(); len(listenAddrs) > 2 {
		return fmt.Errorf("Expected at most 1 listen garlic address, got %v", listenAddrs)
	} else if len(listenAddrs) == 0 {
		// no addr
		return nil
	} else if info, err := t.makePeerInfo(t.ipfsHost.ID(), listenAddrs[1]); err != nil {
		return err
	} else {
		t.peerInfo = info
		return nil
	}
}

func (t *i2pDHT) makePeerInfo(id peer.ID, addr ma.Multiaddr) (*i2pdht.PeerInfo, error) {
	ret := &i2pdht.PeerInfo{ID: id.Pretty()}
	var err error
	if ret.EepServiceID, ret.EepPort, err = defaultAddrFormat.garlicInfo(addr); err != nil {
		return nil, err
	}
	return ret, nil
}

func (t *i2pDHT) connectPeers(ctx context.Context, peers []*i2pdht.PeerInfo, minRequired int) error {
	if len(peers) < minRequired {
		minRequired = len(peers)
	}
	t.debugf("Starting %v peer connections, waiting for at least %v", len(peers), minRequired)
	// Connect to a bunch asynchronously
	peerConnCh := make(chan error, len(peers))
	for _, peer := range peers {
		// There may be an inexplicable race here so I sleep a tad
		// TODO: investigate
		time.Sleep(100 * time.Millisecond)
		go func(peer *i2pdht.PeerInfo) {
			t.debugf("Attempting to connect to peer %v", peer)
			if err := t.connectPeer(ctx, peer); err != nil {
				t.debugf("Failed connecting to peer %v: %v", err)
				peerConnCh <- fmt.Errorf("Peer connection to %v failed: %v", peer, err)
			} else {
				t.debugf("Successfully connected to peer %v", peer)
				peerConnCh <- nil
			}
		}(peer)
	}
	peerErrs := []error{}
	peersConnected := 0
	// Until there is an error or we have enough
	for {
		select {
		case peerErr := <-peerConnCh:
			if peerErr == nil {
				peersConnected++
				if peersConnected >= minRequired {
					return nil
				}
			} else {
				peerErrs = append(peerErrs, peerErr)
				if len(peerErrs) > len(peers)-minRequired {
					return fmt.Errorf("Many failures, unable to get enough peers: %v", peerErrs)
				}
			}
		case <-ctx.Done():
			return fmt.Errorf("Context errored with '%v', peer errors: %v", ctx.Err(), peerErrs)
		}
	}
}

func (t *i2pDHT) connectPeer(ctx context.Context, peerInfo *i2pdht.PeerInfo) error {
	if peer, err := t.addPeer(peerInfo); err != nil {
		return err
	} else {
		return t.ipfsHost.Connect(ctx, *peer)
	}
}

func (t *i2pDHT) addPeer(peerInfo *i2pdht.PeerInfo) (*peerstore.PeerInfo, error) {
	ipfsAddrStr := fmt.Sprintf("%v/ws/ipfs/%v",
		defaultAddrFormat.garlicAddr(peerInfo.EepServiceID, peerInfo.EepPort), peerInfo.ID)
	if ipfsAddr, err := addr.ParseString(ipfsAddrStr); err != nil {
		return nil, err
	} else if peer, err := peerstore.InfoFromP2pAddr(ipfsAddr.Multiaddr()); err != nil {
		return nil, err
	} else {
		t.ipfsHost.Peerstore().AddAddrs(peer.ID, peer.Addrs, peerstore.PermanentAddrTTL)
		return peer, nil
	}
}
