package p2p

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Peer struct {
	Addr         string
	PrivateAddr  string
	ResolvedAddr *net.UDPAddr
}

// Conn needs to be suppled
type PeerManager struct {
	LAddr string
	Peers map[string]*Peer
	mux   sync.RWMutex
	Conn  *net.UDPConn
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

	if _, ok := pm.Peers[addr]; ok {
		delete(pm.Peers, addr)
	}
}

// connects to a peer if peer exists in the list.
// if not exists creates a neww entry in the list and connects to it.
func (pm *PeerManager) Connect(addr string) error {
	if _, ok := pm.Peers[addr]; !ok {
		pm.Peers[addr] = &Peer{
			Addr: addr,
		}
	}

	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	pm.Peers[addr].ResolvedAddr = laddr
	// conn, err := net.Dial("udp", addr)
	// log.Println(pm.LAddr, "Connected to peer:", addr)
	// pm.Peers[addr].Conn = conn

	if pm.Conn == nil {
		return fmt.Errorf("Peer Manager, Connection not assigned")
	}
	return nil
}

// sends a request to the addr to ask for addr's list of peers
func (pm *PeerManager) DiscoverPeers(addr string) error {

	err := pm.Connect(addr)
	if err != nil {
		return err
	}

	if peer, ok := pm.Peers[addr]; ok {
		client := NewPeerClient(*peer, pm.Conn, 2*time.Second)
		list, err := client.GetPeerList()
		if err != nil {
			return err
		}
		pm.AddPeers(list...)

		return nil
	}
	return fmt.Errorf("addr not present in the peer.")
}

func (pm *PeerManager) GetPeerList() []string {
	res := []string{}
	for key := range pm.Peers {
		res = append(res, key)
	}
	return res
}

