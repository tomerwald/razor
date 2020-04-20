package razor

import (
	"../config"
	"../peer_protocol"
	"encoding/binary"
	"log"
	"net"
)

type RazorClient struct {
	Peer              peer_protocol.PeerConnection
	CurrentPiece      uint32
	CurrentBlockIndex uint32
}

func (r *RazorClient) handleChoke() {
	r.Peer.PeerChoking = true
}
func (r *RazorClient) handleUnChoke() {
	r.Peer.PeerChoking = false
}
func (r *RazorClient) handleInterested() {
	r.Peer.PeerInterested = true
}
func (r *RazorClient) handleNotInterested() {
	r.Peer.PeerInterested = false
}
func (r *RazorClient) ChangePiece(haveMsg peer_protocol.Message) {
	r.CurrentPiece = binary.BigEndian.Uint32(haveMsg.Payload)
}
func (r *RazorClient) RequestNextBlock(amount int) {
	var messages []peer_protocol.Message
	for i := 0; i < amount; i++ {
		messages = append(messages, r.Peer.Request(r.CurrentPiece, r.CurrentBlockIndex))
		r.CurrentBlockIndex = r.CurrentBlockIndex + r.Peer.PeerConfig.BlockSize
		if r.CurrentBlockIndex > r.Peer.PeerConfig.PieceSize {
			// next piece
			r.CurrentBlockIndex = 0
			r.CurrentPiece = r.CurrentPiece + 1
		}
	}
	r.Peer.SendMessage(messages...)
}
func (r *RazorClient) MessageCycle() (peer_protocol.Message, error) {
	if m, err := r.Peer.ReceiveMessage(); err == nil {
		switch m.Type {
		case peer_protocol.Bitfield:
			r.Peer.SendMessage(r.Peer.BitField())
		case peer_protocol.Choke:
			r.handleChoke()
		case peer_protocol.Unchoke:
			r.handleUnChoke()
		case peer_protocol.Interested:
			r.handleInterested()
		case peer_protocol.NotInterested:
			r.handleNotInterested()
		case peer_protocol.Have:
			r.ChangePiece(m)
		case peer_protocol.Piece:
			p := ReadPiece(m.Payload)
			ReadCommand(p.Data)
		default:
			log.Printf("Got unknown message type: %d\r\n", m.Type)

		}
		log.Printf("Got message of type %d\r\n", m.Type)
		return m, nil
	} else {
		return m, err
	}
}

func (r *RazorClient) Serve() {
	for r.Peer.Active {
		_, err := r.MessageCycle()
		if err != nil {
			break
		}
		if !r.Peer.PeerChoking {
			log.Println("Requesting command")
			r.RequestNextBlock(1)
		}
	}
}

func NewRazorClient(conn net.Conn, pc config.PeerConfig) RazorClient {
	return RazorClient{peer_protocol.PeerConnection{
		Conn:        conn,
		Active:      false,
		PeerConfig:  pc,
		AmChoking:   false,
		PeerChoking: true,
	},
		0,
		0,
	}
}
