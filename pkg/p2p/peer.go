package p2p

import (
	"log"
	"net"
	"sync"
	"time"
)

type Peer struct {
	Addr string
	Conn net.Conn
}

type PeerManager struct {
	LAddr string
	Peers map[string]*Peer
	mux   sync.RWMutex
}

func NewPeerManager(laddr string) *PeerManager {
	return &PeerManager{
		LAddr: laddr,
		Peers: make(map[string]*Peer),
	}
}

func (pm *PeerManager) AddPeers(peers ...string) {
	pm.mux.Lock()
	defer pm.mux.Unlock()

	for _, addr := range peers {
		if _, ok := pm.Peers[addr]; !ok {
			pm.Peers[addr] = &Peer{Addr: addr}
		}
	}
}

func (pm *PeerManager) RemovePeer(addr string) {
	pm.mux.Lock()
	defer pm.mux.Unlock()

	if peer, ok := pm.Peers[addr]; ok {
		peer.Conn.Close()
		delete(pm.Peers, addr)
	}
}

// connects to a peer if peer exists in the list.
// if not exists creates a neww entry in the list and connects to it.
func (pm *PeerManager) Connect(addr string) (net.Conn, error) {
	if _, ok := pm.Peers[addr]; !ok {
		pm.Peers[addr] = &Peer{
			Addr: addr,
			Conn: nil,
		}
	}
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}
	log.Println(pm.LAddr, "Connected to peer:", addr)
	pm.Peers[addr].Conn = conn
	return conn, nil
}

// sends a request to the addr to ask for addr's list of peers
func (pm *PeerManager) DiscoverPeers(addr string) error {
	conn, err := pm.Connect(addr)
	if err != nil {
		return err
	}

	client := NewPeerClient(conn, 2*time.Second)
	list, err := client.GetPeerList()
	if err != nil {
		return err
	}
	pm.AddPeers(list...)

	return nil
}

func (pm *PeerManager) GetPeerList() []string {
	res := []string{}
	for key := range pm.Peers {
		res = append(res, key)
	}
	return res
}
