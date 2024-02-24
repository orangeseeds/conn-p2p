package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/orangeseeds/udp-holepunching/p2p"
)

var msgRecv chan bool = make(chan bool)
var msgRecv1 chan bool = make(chan bool)

func runNode(laddr string, relayAddr string) {

	flag.Parse()
	node := p2p.NewNode(laddr)
	err := node.Listen()
	if err != nil {
		log.Fatal(err)
	}

	err = handshake(node, relayAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		var msg p2p.Message

		_, addr, err := node.ReadMsg(&msg)
		if err != nil {
			log.Println("Error reading msg", err)
			continue
		}

		log.Printf("recv %v from %v\n", msg.Type.String(), addr.String())
		switch msg.Type {
		case p2p.CONN_FOR:
			err = handleCONN_FOR(node, msg, addr)
			if err != nil {
				log.Println("Error handing conn_for: ", err)
			}
		case p2p.ACPT_FOR:
			err = handleACPT_FOR(node, msg, addr)
			if err != nil {
				log.Println("Error handing acpt_for: ", err)
			}
		case p2p.INIT_PUNCH:
			err = handleINIT_PUNCH(node, msg, addr)
			if err != nil {
				log.Println("Error handing acpt_for: ", err)
			}
		case p2p.MSG:
			log.Println("Test", string(msg.Payload), "from", addr.String())
			// msgRecv <- true
		}
	}
}

func handshake(node *p2p.Node, relayAddr string) error {
	toAddr, err := net.ResolveUDPAddr("udp", relayAddr)
	if err != nil {
		return err
	}

	_, err = node.WriteTo(p2p.Message{
		Type: p2p.SYNC,
		From: node.PublicAddr,
	}, toAddr)
	if err != nil {
		return err
	}

	var msg p2p.Message
	_, _, err = node.ReadMsg(&msg)
	if err != nil {
		return err
	}
	node.PublicAddr = string(msg.Payload)
	log.Println("my addr:", node.PublicAddr)

	err = node.PeerManager.DiscoverPeers(relayAddr)
	if err != nil {
		return err
	}

	val := ""
	fmt.Println("Enter Value: ")
	fmt.Scanf("%s", &val)
	if val == "" {
		log.Println("Skipping")
		return nil
	}
	log.Println("Selected: ", val)
	connMsg := p2p.Message{
		Type: p2p.CONN,
		From: node.PublicAddr,
	}
	timestamp := time.Now().UnixNano()
	connMsg.InjectPayload(p2p.ConnPayload{
		Addr:   val,
		SentAt: timestamp,
	})

	_, err = node.WriteTo(connMsg, toAddr)
	if err != nil {
		return err
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
		From: node.PublicAddr,
	}

	reply.InjectPayload(p2p.ConnPayload{
		Addr:   msg.From,
		SentAt: connPayload.SentAt,
	})

	_, err = node.WriteTo(reply, addr)
	if err != nil {
		return err
	}
	return nil
}

func handleACPT_FOR(node *p2p.Node, msg p2p.Message, addr net.Addr) error {
	var connPayload p2p.ConnPayload
	err := msg.DecodeConnPayload(&connPayload)
	if err != nil {
		return err
	}

	rtt := time.Now().UnixNano() - connPayload.SentAt

	toPeer := p2p.Message{
		Type: p2p.INIT_PUNCH,
		From: node.PublicAddr,
	}
	toPeer.InjectPayload(p2p.ConnPayload{
		Addr:   msg.From, // change this
		SentAt: 0,
	})

	_, err = node.WriteTo(toPeer, addr)
	if err != nil {
		return err
	}

	log.Println("sending INIT_PUNCH to", addr)
	rAddr, err := net.ResolveUDPAddr("udp", msg.From)
	if err != nil {
		return err
	}
	log.Println("sending MSG after t/2", msg.From)

	time.AfterFunc(time.Duration(rtt/2), func() {
		log.Println("sent at", time.Now().UnixNano())
		_, err := node.WriteTo(p2p.Message{
			Type:    p2p.MSG,
			From:    node.PublicAddr,
			Payload: []byte("Apple"),
		}, rAddr)
		if err != nil {
			log.Println(err)
			return
		}
	})
	return nil
}

func handleINIT_PUNCH(node *p2p.Node, msg p2p.Message, addr net.Addr) error {
	rAddr, err := net.ResolveUDPAddr("udp", msg.From)
	if err != nil {
		return err
	}

	log.Println("sent at", time.Now().UnixNano())
	_, err = node.WriteTo(p2p.Message{
		Type:    p2p.MSG,
		From:    node.PublicAddr,
		Payload: []byte("Apple"),
	}, rAddr)
	if err != nil {
		return err
	}
	return nil
}
