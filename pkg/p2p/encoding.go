package p2p

import (
	"io"
)

type MsgDecoder struct{}

// Checks if the message is a message or a stream.
// If it is a message
func (m *MsgDecoder) Decode(r io.Reader, msg *Message) error {

	peekHeader := make([]byte, 1)
	if _, err := r.Read(peekHeader); err != nil {
		return err
	}

	msg.IsStream = peekHeader[0] == 0x2

	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	msg.Payload = buf[:n]
	return nil
}
