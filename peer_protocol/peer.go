package peer_protocol

import (
	"../config"
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type PeerConnection struct {
	Conn         net.Conn
	RemotePeerID []byte
	Active       bool
	PeerConfig   config.PeerConfig
	AmChoking    bool
	PeerChoking  bool
}

func (p *PeerConnection) readSocket(len uint32, timeout int) ([]byte, error) {
	buf := make([]byte, len)
	err := p.Conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(timeout)))
	if err == nil {
		_, err := p.Conn.Read(buf)
		return buf, err
	} else {
		return buf, err
	}
}

func (p *PeerConnection) Disconnect() {
	err := p.Conn.Close()
	if err != nil {
		fmt.Printf("failed closing connection with client %s", p.Conn.RemoteAddr().String())
	} else {
		fmt.Printf("Client disconnected %s", p.Conn.RemoteAddr().String())

	}
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
	buf = append(buf, p.PeerConfig.PeerID...)
	buf = append(buf, p.PeerConfig.InfoHash...)
	_, err := p.Conn.Write(buf)
	if err != nil {
		fmt.Printf("Error while handshaking %s", err.Error())
	} else {
	}
}

func (p *PeerConnection) IsHandshakeValid(protocol string, hash []byte) bool {
	if protocol != protocol {
		return false
	} else if bytes.Compare(hash, p.PeerConfig.InfoHash) == 0 {
		return false
	} else {
		return true
	}
}

func (p *PeerConnection) ReceiveHandshake() bool {
	fmt.Println("Handshaking")
	buf, err := p.readSocket(68, 10)
	if err != nil {
		fmt.Println("Failed to read bytes from socket")
	}
	proto := string(buf[:20])
	p.RemotePeerID = buf[27:47]
	InfoHash := buf[47:67]
	if p.IsHandshakeValid(proto, InfoHash) {
		fmt.Printf("Peer id: %s\r\n", string(p.RemotePeerID))
		p.Active = true
		return true
	} else {
		return false
	}
}

func (p *PeerConnection) ReceiveLength() (uint32, error) {
	lengthBuffer, err := p.readSocket(4, p.PeerConfig.IdleTimeout)
	messageLength := binary.BigEndian.Uint32(lengthBuffer)
	if err == nil {
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
	if messageLength == 0 {
		return Message{Type: KeepAlive}, nil
	}
	if messageBuffer, err := p.readSocket(messageLength, 5); err == nil {
		return readMessage(messageBuffer), err
	} else {
		return Message{}, err
	}
}

func (p *PeerConnection) sendMessage(message Message) {
	if _, err := p.Conn.Write(message.Buffer()); err != nil {
		fmt.Println("Failed sending message")
		p.Active = false
	}
}

func (p *PeerConnection) Choke() {
	m := Message{
		Type: Choke,
	}
	p.sendMessage(m)
}

func (p *PeerConnection) UnChoke() {
	m := Message{
		Type: Unchoke,
	}
	p.sendMessage(m)
}

func (p *PeerConnection) BitField() {
	payload := GenerateRandomBytes(p.PeerConfig.PieceCount)
	m := Message{
		Type:    Bitfield,
		Payload: payload,
	}
	p.sendMessage(m)
}

func (p *PeerConnection) PerformHandshake() bool {
	if p.ReceiveHandshake() {
		p.SendHandshake()
		return true
	} else {
		return false
	}
}

func GenerateRandomBytes(count int) []byte {
	buf := make([]byte, count)
	rand.Read(buf)
	return buf
}
