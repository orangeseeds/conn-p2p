package p2p

import (
	"log"
	"net"
	"sync"
)

type HandlerFunc func(conn net.PacketConn, from net.Addr, msg Message)

type Node struct {
	LAddr       string
	PeerManager *PeerManager
	Listener    net.PacketConn
	Handler     HandlerFunc
	wg          sync.WaitGroup
}

func NewNode(addr string) *Node {
	return &Node{
		LAddr:       addr,
		PeerManager: NewPeerManager(addr),
		Listener:    nil,
		Handler:     nil,
	}
}

func (n *Node) Listen() error {
	n.wg.Add(1)
	conn, err := net.ListenPacket("udp", n.LAddr)
	if err != nil {
		return err
	}
	log.Println("Node is listening at: ", n.LAddr)
	n.Listener = conn

	go func() {
		for {
			buff := make([]byte, 100)
			read, addr, err := conn.ReadFrom(buff)
			if err != nil {
				log.Println("Error reading msg", err)
				continue
			}
			var msg Message
			DecodeMsg(buff[:read], &msg)
			if n.Handler != nil {
				n.Handler(conn, addr, msg)
			}
		}
	}()
	n.wg.Wait()
	return nil
}

func (n *Node) StopListening() {
	if n.Listener != nil {
		n.Listener.Close()
		n.wg.Done()
	}
}

