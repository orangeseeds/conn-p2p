package p2p

import (
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
