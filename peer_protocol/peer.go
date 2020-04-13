package peer_protocol

import (
	"encoding/binary"
	"fmt"
	"net"
)

type PeerConnection struct {
	Conn     net.Conn
	PeerID   []byte
	InfoHash []byte
}

func (p PeerConnection) CreateMessageBuffer(mt byte, payload []byte) []byte {
	LengthField := make([]byte, 4)
	binary.BigEndian.PutUint32(LengthField, uint32(len(payload)+1))
	ActionField := []byte{mt}
	PrefixBuffer := append(ActionField, LengthField...)
	buffer := append(PrefixBuffer, payload...)
	return buffer
}

func (p PeerConnection) Handshake() {
	fmt.Println("Handshaking")
	protocol_type := []byte(Protocol)          // const
	reserved := []byte{0, 0, 0, 0, 0, 0, 0, 0} // all zero 8 bytes
	buf := append(protocol_type, reserved...)
	buf = append(buf, p.PeerID...)
	buf = append(buf, p.InfoHash...)
	l, err := p.Conn.Write(buf)
	if err != nil{}
	fmt.Println(l)
}
