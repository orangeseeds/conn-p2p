package p2p

import (
	"net"
	"testing"
	"time"
)

const (
	addr1 string = "127.0.0.1:1111"
	addr2        = "127.0.0.1:1112"
	addr3        = "127.0.0.1:1113"
)

func TestRelay(t *testing.T) {
	n1 := NewNode(addr1)
	n2 := NewNode(addr2)
	relay := NewNode(addr3)

	err := relay.Listen()
	if err != nil {
		t.Fatal("Listen error:", err)
	}
	err = n1.Listen()
	if err != nil {
		t.Fatal("Listen error:", err)
	}
	err = n2.Listen()
	if err != nil {
		t.Fatal("Listen error:", err)
	}

	// Write CONN Message from n2 to relay
	wMsg := Message{Type: CONN, From: n2.LAddr}
	wMsg.InjectPayload(ConnPayload{
		Addr:   addr2,
		SentAt: time.Now().Unix(),
	})
	_, err = n1.WriteTo(wMsg, relay.ResolvedAddr)
	if err != nil {
		t.Fatal("Writing error: ", err)
	}

	// Read message from Relay
	var rMsg Message
	_, _, err = relay.ReadMsg(&rMsg)
	if err != nil {
		t.Fatal("Read Relay error:", err)
	}
	if rMsg.Type != CONN {
		t.Fatal("Error message type not CONN")
	}
	var connPayload ConnPayload
	err = rMsg.DecodeConnPayload(&connPayload)
	if err != nil {
		t.Fatal("Error decoding:", err)
	}

	toAddr, err := net.ResolveUDPAddr("udp", connPayload.Addr)
	if err != nil {
		t.Fatal("Resolve error:", err)
	}

	// Forward CONN message to n2
	relay.WriteTo(Message{
		Type:    CONN_FOR,
		From:    rMsg.From,
		Payload: rMsg.Payload,
	}, toAddr)

	var n2Msg Message
	_, _, err = n2.ReadMsg(&n2Msg)
	if err != nil {
		t.Fatal("Read error:", err)
	}

	var cPayload ConnPayload
	err = n2Msg.DecodeConnPayload(&cPayload)
	if err != nil {
		t.Fatal("Error decoding:", err)
	}

	// Send ACPT message to n1-relay
	aMsg := Message{Type: ACPT, From: n1.LAddr}
	wMsg.InjectPayload(ConnPayload{
		Addr:   addr2,
		SentAt: time.Now().Unix(),
	})
	_, err = n1.WriteTo(aMsg, relay.ResolvedAddr)
	if err != nil {
		t.Fatal("Writing error: ", err)
	}

	// Read message from Relay
	var raMsg Message
	_, _, err = relay.ReadMsg(&raMsg)
	if err != nil {
		t.Fatal("Read Relay error:", err)
	}
	if raMsg.Type != ACPT {
		t.Fatal("Error message type not CONN")
	}
	var connPayload1 ConnPayload
	err = raMsg.DecodeConnPayload(&connPayload1)
	if err != nil {
		t.Fatal("Error decoding:", err)
	}

	toAddr, err = net.ResolveUDPAddr("udp", connPayload1.Addr)
	if err != nil {
		t.Fatal("Resolve error:", err)
	}
	relay.WriteTo(Message{
		Type:    ACPT_FOR,
		From:    raMsg.From,
		Payload: raMsg.Payload,
	}, toAddr)

	var rn2Msg Message
	_, _, err = n1.ReadMsg(&rn2Msg)
	if err != nil {
		t.Fatal("Read error:", err)
	}

	var rcPayload ConnPayload
	err = rn2Msg.DecodeConnPayload(&rcPayload)
	if err != nil {
		t.Fatal("Error decoding:", err)
	}
	t.Log(rcPayload)

}
