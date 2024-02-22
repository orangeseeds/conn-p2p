package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type MsgType int

const (
	LIST_REQ   MsgType = 0x1 // Requests peers neighbour list
	LIST       MsgType = 0x2 // neighbour list response
	MSG        MsgType = 0x3 // regular string message
	CONN       MsgType = 0x4 // connection request sent to the relay
	CONN_FOR   MsgType = 0x5 // connection request forwarded to the destination
	ACPT       MsgType = 0x6 // connection accept from the destination to the relay
	ACPT_FOR   MsgType = 0x7 // connection accept forwarded to the original sender
	SYNC       MsgType = 0x8 // sync to connect to relay and remain connected
	SYNC_CLOSE MsgType = 0x9
	SYNC_REP   MsgType = 0xA
)

func (m MsgType) String() string {
	switch m {
	case LIST_REQ:
		return "LIST_REQ"
	case LIST:
		return "LIST"
	case MSG:
		return "MSG"
	case CONN:
		return "CONN"
	case CONN_FOR:
		return "CONN_FOR"
	case ACPT:
		return "ACPT"
	case ACPT_FOR:
		return "ACPT_FOR"
	case SYNC:
		return "SYNC"
	case SYNC_CLOSE:
		return "SYNC_CLOSE"
	case SYNC_REP:
		return "SYNC_REP"
	}
	return "NONE"
}
func (m MsgType) Value() int {
	return int(m)
}

type ConnPayload struct {
	Addr   string
	SentAt int64
}

type Message struct {
	Type    MsgType
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
