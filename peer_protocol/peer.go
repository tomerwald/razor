package peer_protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type PeerConnection struct {
	Conn         net.Conn
	PeerID       []byte
	InfoHash     []byte
	remotePeerID []byte
	PieceCount   int
	Active       bool
}

func (p *PeerConnection) createMessageBuffer(mt byte, payload []byte) []byte {
	LengthField := make([]byte, 4)
	binary.BigEndian.PutUint32(LengthField, uint32(len(payload)+1))
	ActionField := []byte{mt}
	PrefixBuffer := append(LengthField, ActionField...)
	buffer := append(PrefixBuffer, payload...)
	return buffer
}

func (p *PeerConnection) SendHandshake() {
	fmt.Println("Handshaking")
	ProtocolType := []byte(Protocol)           // const
	reserved := []byte{0, 0, 0, 0, 0, 0, 0, 0} // all zero 8 bytes
	buf := append(ProtocolType, reserved...)
	buf = append(buf, p.PeerID...)
	buf = append(buf, p.InfoHash...)
	_, err := p.Conn.Write(buf)
	if err != nil {
		fmt.Printf("Error while handshaking %s", err.Error())
	} else {
	}
}

func (p *PeerConnection) IsHandshakeValid(protocol string, hash []byte) bool {
	if protocol != protocol {
		return false
	} else if bytes.Compare(hash, p.InfoHash) == 0 {
		return false
	} else {
		return true
	}
}

func (p *PeerConnection) ReceiveHandshake() bool {
	fmt.Println("Handshaking")
	buf := make([]byte, 68)
	err := p.Conn.SetReadDeadline(time.Now().Local().Add(time.Second + time.Duration(5)))
	if err == nil {
		_, err := p.Conn.Read(buf)
		if err != nil {
			fmt.Println("Failed to read bytes from socket")
		}
	}
	proto := string(buf[:20])
	p.remotePeerID = buf[27:47]
	InfoHash := buf[47:67]
	if p.IsHandshakeValid(proto, InfoHash) {
		fmt.Printf("Peer id: %s\r\n", string(p.remotePeerID))
		p.Active = true
		return true
	} else {
		return false
	}
}

func (p *PeerConnection) ReceiveLength() (uint32, error) {
	lengthBuffer := make([]byte, 4)
	l, err := p.Conn.Read(lengthBuffer)
	messageLength := binary.BigEndian.Uint32(lengthBuffer)
	if l > 0 && err == nil {
		return messageLength, nil
	} else {
		return 0, err
	}
}

func (p *PeerConnection) ReceiveMessage() (Message, error) {
	messageLength, lerr := p.ReceiveLength()
	if lerr != nil {
		return Message{}, lerr
	}
	messageBuffer := make([]byte, messageLength)
	_, err := p.Conn.Read(messageBuffer)
	if err == nil {
		return readMessage(messageBuffer), err
	} else {
		return Message{}, err
	}
}

func (p *PeerConnection) Choke() error {
	m := Message{
		Type:    Choke,
		Payload: nil,
	}
	_, err := p.Conn.Write(m.Buffer())
	return err
}

func (p *PeerConnection) BitField() error {
	payload := GenerateRandomBytes(p.PieceCount)
	m := Message{
		Type:    Bitfield,
		Payload: payload,
	}
	_, err := p.Conn.Write(m.Buffer())
	return err
}

func (p *PeerConnection) PerformHandshake() bool {
	if p.ReceiveHandshake() {
		p.SendHandshake()
		return true
	} else {
		return false
	}
}

func (p *PeerConnection) MessageCycle() {
	m, err := p.ReceiveMessage()
	if err == nil {
		switch m.Type {
		case Bitfield:
			p.BitField()
		default:
			fmt.Printf("Got unknown message type: %d", m.Type)

		}
		fmt.Printf("Got message of type %d", m.Type)
	}
}

func GenerateRandomBytes(count int) []byte {
	buf := make([]byte, count)
	rand.Read(buf)
	return buf
}
