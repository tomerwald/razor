package commands

import (
	"encoding/json"
	"net"
)

type Tunnel struct {
	con     net.Conn
	Timeout int
}

func (t *Tunnel) Close() error {
	return t.con.Close()
}

func NewTunnel(com *Command) (Tunnel, error) {
	var stc StartTunnelCommand
	err := json.Unmarshal(com.Payload, &stc)
	if err != nil {
		return Tunnel{}, err
	}
	con, err := net.Dial("tcp", stc.RemoteAddress)
	return Tunnel{
		con:     con,
		Timeout: stc.timeout,
	}, err
}

type StartTunnelCommand struct {
	RemoteAddress string
	timeout       int
}
