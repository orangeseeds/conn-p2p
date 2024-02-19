package p2p

import (
	"bytes"
	"encoding/gob"
)

const (
	LIST_REQ int = 0x1
	LIST         = 0x2
	MSG          = 0x3
)

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
