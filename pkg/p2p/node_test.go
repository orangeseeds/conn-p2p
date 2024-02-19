package p2p

import (
	"testing"
	"time"
)

func runNode(t *testing.T, laddr string, raddr string) {
	n := NewNode(laddr)
	server, err := n.Listen()
	if err != nil {
		t.Fatal("node listening failed:", err)
	}

	go func() {
		time.Sleep(200 * time.Millisecond)
		conn, err := n.PeerManager.Connect(raddr)
		if err != nil {
			t.Log("connection failed", err)
            return
		}
		conn.Write([]byte("Hello from " + laddr + " to " + raddr))
	}()

	for {
		buff := make([]byte, 100)
		n, addr, err := server.ReadFrom(buff)
		if err != nil {
			t.Fatal(laddr, "Error reading from", addr)
		}
		t.Log(string(buff[:n]), " sent from ", addr)
	}

}

func TestPeerManager(t *testing.T) {
	go func() {
		runNode(t, "127.0.0.1:1111", "127.0.0.1:1112")
	}()
	runNode(t, "127.0.0.1:1112", "127.0.0.1:1111")
}
