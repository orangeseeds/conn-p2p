package main

import (
	"flag"
	"log"
	"net"

	"github.com/orangeseeds/holepunching/pkg/p2p"
)

func main() {
	laddr := flag.String("laddr", "127.0.0.1:1111", "laddr")
	flag.Parse()
	relay := p2p.NewNode(*laddr)
	err := relay.Listen()
	if err != nil {
		log.Fatal(err)
	}

	for {
		var msg p2p.Message

		_, addr, err := relay.ReadMsg(&msg)
		if err != nil {
			log.Println(err)
			continue
		}

        log.Printf("recv %v from %v\n", msg.Type, addr.String())
		switch msg.Type {
		case p2p.CONN:
			handleCONN(relay, msg)
		case p2p.ACPT:
			handleACPT(relay, msg)
		case p2p.SYNC:
			log.Println("received SYNC from: ", msg.From)
		}
	}
}

func handleCONN(relay *p2p.Node, msg p2p.Message) error {
	var connPayload p2p.ConnPayload
	err := msg.DecodeConnPayload(&connPayload)
	if err != nil {
		return err
	}

	toAddr, err := net.ResolveUDPAddr("udp", connPayload.Addr)
	if err != nil {
		return err
	}

	// Forward CONN message to n2
	relay.WriteTo(p2p.Message{
		Type:    p2p.CONN_FOR,
		From:    msg.From,
		Payload: msg.Payload,
	}, toAddr)
	return nil
}

func handleACPT(relay *p2p.Node, msg p2p.Message) error {
	var connPayload p2p.ConnPayload
	err := msg.DecodeConnPayload(&connPayload)
	if err != nil {
		return err
	}

	toAddr, err := net.ResolveUDPAddr("udp", connPayload.Addr)
	if err != nil {
		return err
	}

	// Forward CONN message to n2
	relay.WriteTo(p2p.Message{
		Type:    p2p.ACPT_FOR,
		From:    msg.From,
		Payload: msg.Payload,
	}, toAddr)
	return nil
}
