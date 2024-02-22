package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

const (
	LIST_REQ   int = 0x1 // Requests peers neighbour list
	LIST           = 0x2 // neighbour list response
	MSG            = 0x3 // regular string message
	CONN           = 0x4 // connection request sent to the relay
	CONN_FOR       = 0x5 // connection request forwarded to the destination
	ACPT           = 0x6 // connection accept from the destination to the relay
	ACPT_FOR       = 0x7 // connection accept forwarded to the original sender
	SYNC           = 0x8 // sync to connect to relay and remain connected
	SYNC_CLOSE     = 0x9
)

type ConnPayload struct {
	Addr   string
	SentAt int64
}

type Message struct {
	Type    int
	From    string
	Payload []byte
}

func (m *Message) InjectPayload(payload any) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(payload)
	m.Payload = buff.Bytes()
}

func (m *Message) DecodeConnPayload(p *ConnPayload) error {
	var buff bytes.Buffer
	buff.Write(m.Payload)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&p)
	if err != nil {
		return fmt.Errorf("Error decoding ConnPayload: %v", err)
	}
	return nil
}

func EncodeMsg(msg Message) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(msg)
	return buff.Bytes()
}

func DecodeMsg(data []byte, msg *Message) error {
	var buff bytes.Buffer
	buff.Write(data)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(msg)
	if err != nil {
		return err
	}
	return nil
}
