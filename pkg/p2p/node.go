package p2p

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type HandlerFunc func(conn net.PacketConn, from net.Addr, msg Message)

type Node struct {
	LAddr        string
	ResolvedAddr *net.UDPAddr
	PeerManager  *PeerManager
	Listener     *net.UDPConn
	Handler      HandlerFunc
	wg           sync.WaitGroup
}

func NewNode(addr string) *Node {
	return &Node{
		LAddr:       addr,
		PeerManager: NewPeerManager(addr),
		Listener:    nil,
		Handler:     nil,
	}
}

// Setup the listener and after settingup the listerner,
// ReadMsg can be used to read messages from the connection
func (n *Node) Listen() error {
	laddr, err := net.ResolveUDPAddr("udp", n.LAddr)
	if err != nil {
		return err
	}
	n.ResolvedAddr = laddr
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}
	log.Println("Node is listening at: ", n.LAddr)
	n.Listener = conn
    n.PeerManager.Conn = conn

	// for {
	// 	buff := make([]byte, 100)
	// 	read, addr, err := conn.ReadFrom(buff)
	// 	if err != nil {
	// 		log.Println("Error reading msg", err)
	// 		continue
	// 	}
	// 	var msg Message
	// 	err = DecodeMsg(buff[:read], &msg)
	// 	if err != nil {
	// 		log.Println("Error decoding msg", err)
	// 		continue
	// 	}
	//
	// 	if n.Handler == nil {
	// 		return fmt.Errorf("Handler is not assigned.")
	// 	}
	// 	n.Handler(conn, addr, msg)
	// }
	return nil
}

func (n *Node) ReadMsg(msg *Message) (int, net.Addr, error) {
	if n.Listener == nil {
		return 0, nil, fmt.Errorf("Listener not assigned")
	}

	buff := make([]byte, 1024)
	read, addr, err := n.Listener.ReadFrom(buff)
	if err != nil {
		return read, addr, err
	}
	err = DecodeMsg(buff[:read], msg)
	if err != nil {
		return read, addr, fmt.Errorf("Error decoding msg %v", err)
	}
	return read, addr, nil
}

func (n *Node) WriteTo(msg Message, addr net.Addr) (int, error) {
	return n.Listener.WriteTo(EncodeMsg(msg), addr)
}

func (n *Node) StopListening() {
	if n.Listener != nil {
		n.Listener.Close()
	}
}

// func (n *Node) Broadcast(msg Message) {
// 	for _, p := range n.PeerManager.Peers {
// 		if p != nil && p.Conn != nil {
// 			p.Conn.Write(EncodeMsg(msg))
// 		}
// 	}
// }
