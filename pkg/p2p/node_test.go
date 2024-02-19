package p2p

import (
	"net"
	"testing"
	"time"
)

func TestDiscoverPeers(t *testing.T) {

	n1 := NewNode("127.0.0.1:1111")
	n2 := NewNode("127.0.0.1:1112")

	n1.Handler = func(conn net.PacketConn, from net.Addr, msg Message) {
		switch msg.Type {
		case LIST_REQ:
			resp := Message{
				Type: LIST,
				From: n1.LAddr,
			}
			resp.InjectPayload(n1.PeerManager.GetPeerList())
			conn.WriteTo(EncodeMsg(resp), from)
		}
	}
	go func() {
		err := n1.Listen()
		if err != nil {
			t.Log("node listening failed:", err)
			return
		}
	}()

	err := n2.PeerManager.DiscoverPeers(n1.LAddr)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range n2.PeerManager.Peers {
		t.Logf("%v", p.Addr)
	}
}
func TestPeerMsg(t *testing.T) {

	payload_msg := "hello there"
	n1 := NewNode("127.0.0.1:1111")
	n2 := NewNode("127.0.0.1:1112")

	n1.Handler = func(conn net.PacketConn, from net.Addr, msg Message) {
		switch msg.Type {
		case MSG:
			if string(msg.Payload) != payload_msg {
				t.Fatalf("expected %s got %s", payload_msg, msg.Payload)
			}
			// t.Logf("Message: %v from: %v to %v", string(msg.Payload), msg.From, conn.LocalAddr().String())
		}
	}
	go func() {
		err := n1.Listen()
		if err != nil {
			t.Log("node listening failed:", err)
			return
		}
	}()

	conn, err := n2.PeerManager.Connect(n1.LAddr)
	if err != nil {
		t.Fatal(err)
	}

	msg := Message{
		Type: MSG,
		From: n2.LAddr,
        Payload: []byte(payload_msg),
	}
	conn.Write(EncodeMsg(msg))
	time.Sleep(200 * time.Millisecond)
}

// func runNode(t *testing.T, laddr string, raddr string) {
// 	n := NewNode(laddr)
// 	server, err := n.Listen()
// 	if err != nil {
// 		t.Fatal("node listening failed:", err)
// 	}
//
// 	go func() {
// 		time.Sleep(200 * time.Millisecond)
// 		conn, err := n.PeerManager.Connect(raddr)
// 		if err != nil {
// 			t.Log("connection failed", err)
// 			return
// 		}
// 		conn.Write([]byte("Hello from " + laddr + " to " + raddr))
// 	}()
//
// 	for {
// 		buff := make([]byte, 100)
// 		n, addr, err := server.ReadFrom(buff)
// 		if err != nil {
// 			t.Fatal(laddr, "Error reading from", addr)
// 		}
// 		t.Log(string(buff[:n]), " sent from ", addr)
// 	}
//
// }
//
// func TestPeerManager(t *testing.T) {
// 	go func() {
// 		runNode(t, "127.0.0.1:1111", "127.0.0.1:1112")
// 	}()
// 	runNode(t, "127.0.0.1:1112", "127.0.0.1:1111")
// }
