package razor

import (
	"../config"
	"../peer_protocol"
	"encoding/binary"
	"fmt"
	"net"
)

type RazorClient struct {
	Peer         peer_protocol.PeerConnection
	CurrentPiece uint32
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
func (r *RazorClient) MessageCycle() (peer_protocol.Message, error) {
	if m, err := r.Peer.ReceiveMessage(); err == nil {
		switch m.Type {
		case peer_protocol.Bitfield:
			r.Peer.SendMessage(r.Peer.BitField(), r.Peer.Have(23))
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
		default:
			fmt.Printf("Got unknown message type: %d\r\n", m.Type)

		}
		fmt.Printf("Got message of type %d\r\n", m.Type)
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
	}
}

func NewRazorClient(conn net.Conn, pc config.PeerConfig) RazorClient {
	return RazorClient{peer_protocol.PeerConnection{
		Conn:        conn,
		Active:      false,
		PeerConfig:  pc,
		AmChoking:   false,
		PeerChoking: false,
	},
		0,
	}
}
