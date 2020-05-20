package peer_protocol

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"time"

	"../config"
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
	// create a new peer using a connected socket and configuration
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
	// close the connection to the peer
	err := p.Conn.Close()
	if err != nil {
		log.Printf("failed closing connection with client %s", p.Conn.RemoteAddr().String())
	} else {
		log.Printf("Client disconnected %s", p.Conn.RemoteAddr().String())

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
	// send a bittorrent handshare containing the peerID info hash etc.
	log.Println("Handshaking")
	ProtocolType := []byte(Protocol)           // const
	reserved := []byte{0, 0, 0, 0, 0, 0, 0, 0} // all zero 8 bytes
	buf := append(ProtocolType, reserved...)
	buf = append(buf, p.Config.PeerID...)
	buf = append(buf, p.Config.InfoHash...)
	_, err := p.Conn.Write(buf)
	if err != nil {
		log.Printf("Error while handshaking %s", err.Error())
	} else {
	}
}

func (p *PeerConnection) IsHandshakeValid(protocol string, hash []byte) bool {
	// validate the given info hash and protocol are as expected
	return bytes.Compare(hash, p.Config.InfoHash) != 0 && protocol == Protocol
}

func (p *PeerConnection) ReceiveHandshake() bool {
	// Receive a bittorrent handshake with the remote peer id and info hash
	log.Println("Handshaking")
	buf, err := p.readSocket(HandshakeLength, 10)
	if err != nil {
		log.Println("Failed to read bytes from socket")
	}
	proto := string(buf[:ProtocolFieldLength])
	p.RemotePeerID = buf[PeerIDOffset : PeerIDFieldLength+PeerIDOffset-1]
	InfoHash := buf[InfoHashOffset : InfoHashFieldLength+InfoHashOffset-1]
	if p.IsHandshakeValid(proto, InfoHash) {
		log.Printf("Peer id: %s\r\n", string(p.RemotePeerID))
		p.Active = true
		return true
	} else {
		return false
	}
}

func (p *PeerConnection) ReceiveLength() (uint32, error) {
	// receive the length of the following message
	lengthBuffer, err := p.readSocket(4, p.Config.IdleTimeout)
	messageLength := binary.BigEndian.Uint32(lengthBuffer)
	if err == nil {
		return messageLength, nil
	} else {
		return 0, err
	}
}

func (p *PeerConnection) ReceiveMessage() (Message, error) {
	// receive a whole message from peer
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
	// send a message struct to peer
	var messageBuffer []byte
	for _, msg := range messages {
		messageBuffer = append(messageBuffer, msg.Buffer()...)
	}
	if _, err := p.Conn.Write(messageBuffer); err != nil {
		log.Println("Failed sending message")
		p.Active = false
	}
}
func (p *PeerConnection) Choke() Message {
	// create a choke message
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
	// create a random BitField message
	payload := GenerateRandomBytes(p.Config.PieceCount)
	return Message{
		Type:    Bitfield,
		Payload: payload,
	}
}
func (p *PeerConnection) Interested() Message {
	// create an Interested message 
	return Message{
		Type: Interested,
	}
}
func (p *PeerConnection) Have(index uint32) Message {
	// create a message informing the peer you obtain a piece (original protol)
	// tell the peer which piece he should ask for next
	indexField := make([]byte, 4)
	binary.BigEndian.PutUint32(indexField, index)
	return Message{
		Type:    Have,
		Payload: indexField,
	}
}
func (p *PeerConnection) Request(PieceIndex uint32, fromOffset uint32) Message {
	// request a block from a specific piece from the peer
	// used to ask for a command
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
	// wait for peer to send a handshake then respond
	if p.ReceiveHandshake() {
		p.SendHandshake()
		return true
	} else {
		return false
	}
}

func GenerateRandomBytes(count int) []byte {
	// generate a random binary blob of a given length
	buf := make([]byte, count)
	rand.Read(buf)
	return buf
}
