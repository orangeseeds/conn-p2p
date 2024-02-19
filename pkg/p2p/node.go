package p2p

import (
	"log"
	"net"
)

type Node struct {
	LAddr       string
	PeerManager *PeerManager
	Listener    net.PacketConn
}

func NewNode(addr string) *Node {
	return &Node{
		LAddr:       addr,
		PeerManager: NewPeerManager(addr),
		Listener:    nil,
	}
}

func (n *Node) Listen() (net.PacketConn, error) {
	conn, err := net.ListenPacket("udp", n.LAddr)
	if err != nil {
		return nil, err
	}
	log.Println("Node is listening at: ", n.LAddr)
	n.Listener = conn
	return conn, nil
}

func (n *Node) StopListening() {
	if n.Listener != nil {
		n.Listener.Close()
	}
}


