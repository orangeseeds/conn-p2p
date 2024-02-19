package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

// handles all activities related to connection with other peers
type PeerClient struct {
	Conn    net.Conn
	Timeout time.Duration
}

func NewPeerClient(conn net.Conn, timeout time.Duration) *PeerClient {
	return &PeerClient{
		Conn:    conn,
		Timeout: timeout,
	}
}

func (pc *PeerClient) GetPeerList() ([]string, error) {
	_, err := pc.Conn.Write(EncodeMsg(Message{
		Type: LIST_REQ,
		From: pc.Conn.LocalAddr().String(),
	}))
	if err != nil {
		return nil, err
	}

	data := make([]byte, 100)
	n, err := pc.Conn.Read(data)
	if err != nil {
		return nil, err
	}

	var msg Message
	DecodeMsg(data[:n], &msg)

	if msg.Type != LIST {
		return nil, fmt.Errorf("Invalid message type: %v, expected %v", msg.Type, LIST)
	}

	var buff bytes.Buffer
	var peers []string
	buff.Write(msg.Payload)
	dec := gob.NewDecoder(&buff)
	err = dec.Decode(&peers)
	if err != nil {
		return nil, err
	}
	return peers, nil
}

func (pc *PeerClient) SendMsg(data string) error {
	msg := Message{
		Type: MSG,
		From: pc.Conn.LocalAddr().String(),
	}
	msg.InjectPayload(data)
	_, err := pc.Conn.Write(EncodeMsg(msg))
	if err != nil {
		return err
	}
	return nil
}
