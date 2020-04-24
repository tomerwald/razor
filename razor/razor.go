package razor

import (
	"../config"
	"../peer_protocol"
	"./commands"
	"encoding/binary"
	"log"
	"net"
)

type Client struct {
	Peer              peer_protocol.PeerConnection
	CurrentPiece      uint32
	CurrentBlockIndex uint32
	CommandOutput     []byte
}

func (r *Client) handleChoke() {
	r.Peer.PeerChoking = true
}
func (r *Client) handleUnChoke() {
	r.Peer.PeerChoking = false
}
func (r *Client) handleInterested() {
	r.Peer.PeerInterested = true
}
func (r *Client) handleNotInterested() {
	r.Peer.PeerInterested = false
}
func (r *Client) ChangePiece(haveMsg peer_protocol.Message) {
	r.CurrentPiece = binary.BigEndian.Uint32(haveMsg.Payload)
}
func (r *Client) HandleRejection() {
	if r.CurrentBlockIndex == 0 {
		r.CurrentBlockIndex = r.Peer.Config.PieceSize - r.Peer.Config.BlockSize
		r.CurrentPiece = r.CurrentPiece - 1
	} else {
		r.CurrentBlockIndex = r.CurrentBlockIndex - r.Peer.Config.PieceSize
	}
}
func (r *Client) RequestNextBlock(amount int) {
	var messages []peer_protocol.Message
	for i := 0; i < amount; i++ {
		req := peer_protocol.NewRequest(r.CurrentPiece, r.CurrentBlockIndex, r.Peer.Config.BlockSize)
		messages = append(messages, req.Message())
		r.CurrentBlockIndex = r.CurrentBlockIndex + r.Peer.Config.BlockSize
		if r.CurrentBlockIndex > r.Peer.Config.PieceSize {
			// next piece
			r.CurrentBlockIndex = 0
			r.CurrentPiece = r.CurrentPiece + 1
		}
	}
	r.Peer.SendMessage(messages...)
}

func (r *Client) createPieceResponse(req peer_protocol.RequestMessage) []byte {
	var chunk []byte
	if len(r.CommandOutput) > int(req.BlockEnd()) {
		chunk = r.CommandOutput[req.BlockOffset:req.BlockEnd()]
	} else {
		chunk = r.CommandOutput[req.BlockOffset:]
	}
	com := commands.Command{
		Type:    0,
		Payload: chunk,
	}
	paddingSize := int(req.BlockSize+8) - com.BufferLen()
	return append(com.Buffer(), peer_protocol.GenerateRandomBytes(paddingSize)...)
}
func (r *Client) RespondToRequest(req peer_protocol.RequestMessage) {
	chunk := r.createPieceResponse(req)
	m := peer_protocol.PieceMessage{
		PieceIndex:  req.PieceIndex,
		BlockOffset: req.BlockOffset,
		Data:        chunk,
	}
	r.Peer.SendMessage(m.Message())
}

func (r *Client) MessageCycle() (peer_protocol.Message, error) {
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
		case peer_protocol.Reject:
			r.HandleRejection()
		case peer_protocol.Piece:
			p := peer_protocol.ReadPiece(m.Payload)
			r.CommandOutput, _ = commands.ReadCommand(p.Data)
		case peer_protocol.Request:
			req := peer_protocol.RequestFromPayload(m.Payload)
			r.RespondToRequest(req)
		default:
			log.Printf("Got unknown message type: %d\r\n", m.Type)
		}
		log.Printf("Got message of type %d\r\n", m.Type)
		return m, nil
	} else {
		return m, err
	}
}

func (r *Client) Serve() {
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

func NewRazorClient(conn net.Conn, pc config.PeerConfig) Client {
	var out []byte
	return Client{peer_protocol.PeerConnection{
		Conn:        conn,
		Active:      false,
		Config:      pc,
		AmChoking:   false,
		PeerChoking: true,
	},
		0,
		0,
		out,
	}
}
