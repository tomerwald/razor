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
	Conn           net.Conn
	RemotePeerID   []byte
	Active         bool
	Config         *config.PeerConfig
	AmChoking      bool
	PeerChoking    bool
	PeerInterested bool
	bitField       []byte
}

func NewPeer(conn net.Conn, config *config.PeerConfig) PeerConnection {
	return PeerConnection{
		Conn:        conn,
		Config:      config,
		PeerChoking: true,
	}
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
	buf = append(buf, p.Config.PeerID...)
	buf = append(buf, p.Config.InfoHash...)
	_, err := p.Conn.Write(buf)
	if err != nil {
		fmt.Printf("Error while handshaking %s", err.Error())
	} else {
	}
}

func (p *PeerConnection) IsHandshakeValid(protocol string, hash []byte) bool {
	if protocol != protocol {
		return false
	} else if bytes.Compare(hash, p.Config.InfoHash) == 0 {
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
	lengthBuffer, err := p.readSocket(4, p.Config.IdleTimeout)
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

func (p *PeerConnection) SendMessage(messages ...Message) {
	messageBuffer := []byte{}
	for _, msg := range messages {
		messageBuffer = append(messageBuffer, msg.Buffer()...)
	}
	if _, err := p.Conn.Write(messageBuffer); err != nil {
		fmt.Println("Failed sending message")
		p.Active = false
	}
}
func (p *PeerConnection) Choke() Message {
	return Message{
		Type: Choke,
	}
}

func (p *PeerConnection) UnChoke() Message {
	return Message{
		Type: Unchoke,
	}
}

func (p *PeerConnection) BitField() Message {
	payload := GenerateRandomBytes(p.Config.PieceCount)
	return Message{
		Type:    Bitfield,
		Payload: payload,
	}
}
func (p *PeerConnection) Interested() Message {
	return Message{
		Type: Interested,
	}
}
func (p *PeerConnection) Have(index uint32) Message {
	indexField := make([]byte, 4)
	binary.BigEndian.PutUint32(indexField, index)
	return Message{
		Type:    Have,
		Payload: indexField,
	}
}
func (p *PeerConnection) Request(PieceIndex uint32, fromOffset uint32) Message {
	indexField := make([]byte, 12)
	binary.BigEndian.PutUint32(indexField[0:4], PieceIndex)
	binary.BigEndian.PutUint32(indexField[4:8], fromOffset)
	binary.BigEndian.PutUint32(indexField[8:12], p.Config.BlockSize)
	return Message{
		Type:    Request,
		Payload: indexField,
	}
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
