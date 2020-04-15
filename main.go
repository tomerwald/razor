package main

import (
	"./config"
	"./peer_protocol"
	"fmt"
	"net"
)

type PeerManger struct {
	Config config.RazorConifg
}

func (pm PeerManger) handleClient(con net.Conn) {
	fmt.Printf("Handling client %s\r\n", con.RemoteAddr().String())
	p := peer_protocol.PeerConnection{
		PeerID:   pm.Config.PeerID,
		InfoHash: pm.Config.InfoHash,
		Conn:     con,
	}
	if p.PerformHandshake() {
		p.Choke()
	} else {

	}
}

func (pm PeerManger) Start() {
	fmt.Println("Getting connections\r\n")
	l, err := net.Listen("tcp", pm.Config.Address)
	if err != nil {
		fmt.Println("Failed listening for connections")
		return
	} else {
		for i := 0; i < 3; i++ {
			conn, err := l.Accept()
			if err != nil {
				continue
			}
			go pm.handleClient(conn)
		}
	}
}
func main() {
	InfoHash := []byte{0xd5, 0x5f, 0x1e, 0x84, 0x0f, 0x1b, 0xd6,
		0x57, 0x6e, 0xad, 0x67, 0xa4, 0xd0, 0x4e, 0x5d, 0x6e, 0xa2, 0x94, 0x41, 0x4b}
	PeerID := []byte("-UW109K-LMYpj9A)8X0R")
	c := config.RazorConifg{
		Address:  "127.0.0.1:6888",
		InfoHash: InfoHash,
		PeerID:   PeerID,
	}
	pm := PeerManger{c}
	pm.Start()
}
