package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/orangeseeds/holepunching/pkg/p2p"
)

func main() {
	laddr := flag.String("laddr", "127.0.0.1:1111", "laddr")
	relayAddr := flag.String("relayAddr", "127.0.0.1:1112", "relay addr")
	to := flag.String("raddr", "127.0.0.1:1113", "raddr")

	flag.Parse()
	time.Sleep(5 * time.Second)
	node := p2p.NewNode(*laddr)
	err := node.Listen()
	if err != nil {
		log.Fatal(err)
	}

	err = handshake(node, *relayAddr, *to)
	if err != nil {
		log.Fatal(err)
	}
	for {
		var msg p2p.Message

		_, addr, err := node.ReadMsg(&msg)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("recv %v from %v\n", msg.Type.String(), addr.String())
		switch msg.Type {
		case p2p.CONN_FOR:
			handleCONN_FOR(node, msg, addr)
		case p2p.ACPT_FOR:
			handleACPT_FOR(node, msg, addr)
		case p2p.SYNC_REP:
			node.PublicAddr = string(msg.Payload)
		}
	}
}

func handshake(node *p2p.Node, relayAddr string, to string) error {
	toAddr, err := net.ResolveUDPAddr("udp", relayAddr)
	if err != nil {
		return err
	}

	_, err = node.WriteTo(p2p.Message{
		Type: p2p.SYNC,
		From: node.LAddr,
	}, toAddr)
	if err != nil {
		return err
	}

	err = node.PeerManager.DiscoverPeers(relayAddr)
	if err != nil {
		return err
	}

	for _, val := range node.PeerManager.Peers {
		connMsg := p2p.Message{
			Type: p2p.CONN,
			From: node.LAddr,
		}
		connMsg.InjectPayload(p2p.ConnPayload{
			Addr:   val.Addr,
			SentAt: time.Now().Unix(),
		})
		_, err = node.WriteTo(connMsg, toAddr)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleCONN_FOR(node *p2p.Node, msg p2p.Message, addr net.Addr) error {
	var connPayload p2p.ConnPayload
	err := msg.DecodeConnPayload(&connPayload)
	if err != nil {
		return err
	}
	reply := p2p.Message{
		Type: p2p.ACPT,
		From: node.LAddr,
	}

	reply.InjectPayload(p2p.ConnPayload{
		Addr:   msg.From,
		SentAt: time.Now().Unix(),
	})

	roundTime := time.Now().Unix() - connPayload.SentAt

	go func() {
		<-time.After(time.Duration(roundTime))
		log.Println("Sent payload to", connPayload.Addr)
	}()

	node.WriteTo(reply, addr)
	return nil
}

func handleACPT_FOR(node *p2p.Node, msg p2p.Message, addr net.Addr) error {
	var connPayload p2p.ConnPayload
	err := msg.DecodeConnPayload(&connPayload)
	if err != nil {
		return err
	}

	log.Println("Sent payload to", addr)
	return nil

}
