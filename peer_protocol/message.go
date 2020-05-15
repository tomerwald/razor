package peer_protocol

import "encoding/binary"

type Message struct {
	// a representation of a bittorrent message
	Type    byte
	Payload []byte
}

func (m Message) Buffer() []byte {
	// Parse the message to binary form
	// format: <length uint32><action-code byte><payload bytes>
	LengthField := make([]byte, 4)
	binary.BigEndian.PutUint32(LengthField, uint32(len(m.Payload)+1))
	ActionField := []byte{m.Type}
	PrefixBuffer := append(LengthField, ActionField...)
	buffer := append(PrefixBuffer, m.Payload...)
	return buffer
}

func readMessage(buf []byte) Message {
	MessageType := buf[0]
	payload := buf[1:]
	return Message{
		Type:    MessageType,
		Payload: payload,
	}
}
