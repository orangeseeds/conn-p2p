package relay

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"net"
	"testing"
)

func TestMsgEncode(t *testing.T) {
	var buffer bytes.Buffer
	src := net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 1111,
	}
	dest := net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 1112,
	}

	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(Message{src, dest, SYNC})
	if err != nil {
		t.Fatal("Error encoding Message:", err)
	}

	dec := gob.NewDecoder(&buffer)
	var m Message
	err = dec.Decode(&m)
	if err != nil {
		t.Fatal("Error decoding Message:", err)
	}
}

func TestRelay(t *testing.T) {
	relay := Relay{
		Port:    ":8080",
		Records: make([]NATRecord, 10),
	}
	go relay.ListenUDP()

	conn, err := net.Dial("udp", "localhost:8080")
	if err != nil {

		t.Fatal(err)
	}

	var buffer bytes.Buffer
	src := net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 1111,
	}
	dest := net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 1112,
	}

	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(Message{src, dest, CONN})
	if err != nil {
		t.Fatal("Error encoding Message:", err)
	}
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	bt := make([]byte, 100)
	_, err = bufio.NewReader(conn).Read(bt)
}

func TestRunRelay(t *testing.T) {
	relay := Relay{
		Port:    ":8080",
		Records: make([]NATRecord, 10),
	}
		relay.ListenUDP()
}
