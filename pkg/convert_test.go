package main

import (
	"net"
	"testing"

	"github.com/orangeseeds/holepunching/pkg/p2p"
)

const addr = "127.0.0.1"

func TestMain(t *testing.T) {
	laddr, err := net.ResolveUDPAddr("udp", addr+":8000")
	if err != nil {
		t.Fatal("Resolve error:", err)
	}
	raddr, err := net.ResolveUDPAddr("udp", addr+":8080")
	if err != nil {
		t.Fatal("Resolve error:", err)
	}

	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		t.Fatal("Dial error:", err)
	}

	_, err = conn.Write([]byte("hello world!"))
	if err != nil {
		t.Fatal("Write error:", err)
	}

	msg := p2p.Message{}
	decoder := p2p.MsgDecoder{}
	decoder.Decode(conn, &msg)
	if err != nil {
		t.Fatal("Payload Error:", err)
	}
	t.Log("Payload", msg)
}
