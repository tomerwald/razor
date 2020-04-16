package razor

import (
	"../config"
	"../peer_protocol"
	"fmt"
	"net"
)

type RazorClient struct {
	Peer peer_protocol.PeerConnection
}

func (r *RazorClient) handleChoke() {
	r.Peer.PeerChoking = true
}
func (r *RazorClient) handleUnChoke() {
	r.Peer.PeerChoking = false
}
func (r *RazorClient) MessageCycle() error {
	if m, err := r.Peer.ReceiveMessage(); err == nil {
		switch m.Type {
		case peer_protocol.Bitfield:
			r.Peer.BitField()
		case peer_protocol.Choke:
			r.handleChoke()
		case peer_protocol.Unchoke:
			r.handleUnChoke()

		default:
			fmt.Printf("Got unknown message type: %d\r\n", m.Type)

		}
		fmt.Printf("Got message of type %d\r\n", m.Type)
		return nil
	} else {
		return err
	}
}

func (r *RazorClient) Serve(){
	for r.Peer.Active {
		err := r.MessageCycle()
		if err != nil {
			break
		}
	}
}

func NewRazorClient(conn net.Conn, pc config.PeerConfig) RazorClient {
	r := RazorClient{peer_protocol.PeerConnection{
		Conn:        conn,
		Active:      false,
		PeerConfig:  pc,
		AmChoking:   false,
		PeerChoking: false,
	},
	}
	return r
}
