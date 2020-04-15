package peer_protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type PeerConnection struct {
	Conn         net.Conn
	PeerID       []byte
	InfoHash     []byte
	remotePeerID []byte
}

func (p PeerConnection) createMessageBuffer(mt byte, payload []byte) []byte {
	LengthField := make([]byte, 4)
	binary.BigEndian.PutUint32(LengthField, uint32(len(payload)+1))
	ActionField := []byte{mt}
	PrefixBuffer := append(LengthField, ActionField...)
	buffer := append(PrefixBuffer, payload...)
	return buffer
}

func (p PeerConnection) SendHandshake() {
	fmt.Println("Handshaking")
	protocol_type := []byte(Protocol)          // const
	reserved := []byte{0, 0, 0, 0, 0, 0, 0, 0} // all zero 8 bytes
	buf := append(protocol_type, reserved...)
	buf = append(buf, p.PeerID...)
	buf = append(buf, p.InfoHash...)
	l, err := p.Conn.Write(buf)
	if err != nil {
	}
	fmt.Println(l)
}

func (p PeerConnection) IsHandshakeValid(protocol string, hash []byte) bool {
	if protocol != protocol {
		return false
	} else if bytes.Compare(hash, p.InfoHash) == 0 {
		return false
	} else {
		return true
	}
}

func (p PeerConnection) ReceiveHandshake() bool {
	fmt.Println("Handshaking")
	buf := make([]byte, 68)
	_, err := p.Conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to read bytes from socket")
	}
	proto := string(buf[:20])
	p.remotePeerID = buf[27:47]
	InfoHash := buf[47:67]
	if p.IsHandshakeValid(proto, InfoHash) {
		fmt.Printf("Peer id: %s\r\n", string(p.remotePeerID))
		return true
	} else {
		return false
	}
}

func (p PeerConnection) Choke() error {
	var payload []byte
	buf := p.createMessageBuffer(Choke, payload)
	_, err := p.Conn.Write(buf)
	return err
}

func (p PeerConnection) PerformHandshake() bool{
	if p.ReceiveHandshake() {
		p.SendHandshake()
		return true
	} else {
		return false
	}
}