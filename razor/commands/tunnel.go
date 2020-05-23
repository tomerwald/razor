package commands

import (
	"encoding/hex"
	"encoding/json"
	"net"
	"time"
)

type Tunnel struct {
	con     net.Conn
	Timeout int
}

func (t *Tunnel) Close() error {
	return t.con.Close()
}
func (t *Tunnel) Send(com *Command) []byte {
	var ts TunnelSend
	err := json.Unmarshal(com.Payload, &ts)
	if err != nil {
		return NewTunnelResponse(true, err.Error()).toBuffer()
	} else {
		binData, _ := hex.DecodeString(ts.Buffer)
		_, err = t.con.Write(binData)
		return NewTunnelResponse(false, "Sent").toBuffer()
	}
}
func (t *Tunnel) Recv(com *Command) []byte {
	var tr TunnelRecv
	err := json.Unmarshal(com.Payload, &tr)
	if err != nil {
		return NewTunnelResponse(true, err.Error()).toBuffer()
	}
	out := make([]byte, tr.ByteCount)
	_ = t.con.SetReadDeadline(time.Now().Add(time.Second * time.Duration(t.Timeout)))
	_, err = t.con.Read(out)
	if err != nil {
		return NewTunnelResponse(true, err.Error()).toBuffer()
	}
	return NewTunnelResponse(false, hex.EncodeToString(out)).toBuffer()
}

func NewTunnel(com *Command) (Tunnel, []byte) {
	var stc StartTunnelCommand
	err := json.Unmarshal(com.Payload, &stc)
	if err != nil {
		return Tunnel{}, NewTunnelResponse(true, err.Error()).toBuffer()
	}
	con, err := net.Dial("tcp", stc.RemoteAddress)
	return Tunnel{
		con:     con,
		Timeout: stc.Timeout,
	}, NewTunnelResponse(false, "Tunnel is up").toBuffer()
}

type StartTunnelCommand struct {
	RemoteAddress string
	Timeout       int
}

type TunnelSend struct {
	Buffer string
}
type TunnelRecv struct {
	ByteCount int
}
type TunnelResponse struct {
	IsError bool
	Content string
}

func (tr *TunnelResponse) toBuffer() []byte {
	out, _ := json.Marshal(tr)
	return out
}
func NewTunnelResponse(isError bool, content string) *TunnelResponse {
	return &TunnelResponse{
		IsError: isError,
		Content: content,
	}
}
