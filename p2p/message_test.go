package p2p

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestMessageEncoding(t *testing.T) {
	msg := Message{
		Type:    LIST_REQ,
		Payload: []byte("Message Payload"),
	}

	encoded := EncodeMsg(msg)
	var decodedMsg Message

	DecodeMsg(encoded, &decodedMsg)
	t.Logf("%v", decodedMsg)
}

func TestConnMessage(t *testing.T) {
	msg := Message{Type: CONN, From: "127.0.0.1:1111"}
	// msg.InjectPayload(ConnPayload{
	// 	Addr:   "127.0.0.1:1112",
	// 	SentAt: time.Now().Unix(),
	// })
	msg.Payload = []byte{
		44, 127, 3, 1, 1, 11, 67, 111, 110, 110, 80, 97, 121, 108, 111, 97, 100, 1, 255, 128, 0, 1, 2, 1, 4, 65, 100, 100, 114, 1, 12, 0, 1, 6, 83, 101, 110, 116, 65, 116, 1, 4, 0, 0, 0, 25, 255, 128, 1, 14, 49, 50, 55, 46, 48, 46, 48, 46, 49, 58, 49, 49, 49, 50, 1, 252, 203, 173, 133, 16, 0,
	}

	encoded := EncodeMsg(msg)

	var dMsg Message
	err := DecodeMsg(encoded, &dMsg)
	if err != nil {
		t.Fatal("Decoding error:", err)
	}

	var connPayload ConnPayload
	var buff bytes.Buffer
	buff.Write(dMsg.Payload)
	dec := gob.NewDecoder(&buff)
	err = dec.Decode(&connPayload)
	if err != nil {
		t.Fatal("Error decoding:", err)
	}
	t.Log(dMsg, connPayload)

}
