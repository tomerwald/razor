package main

import (
	"./config"
	"./peer_protocol"
	"fmt"
	"net"
)

type PeerManger struct {
	ManConfig  config.ManagerConfig
	PeerConfig config.PeerConfig
}

func (pm *PeerManger) handleClient(con net.Conn) {
	fmt.Printf("Handling client %s\r\n", con.RemoteAddr().String())
	p := peer_protocol.PeerConnection{
		Conn:       con,
		PeerConfig: pm.PeerConfig,
	}
	defer con.Close()
	if p.PerformHandshake() {
		for p.Active {
			err := p.MessageCycle()
			if err != nil {
				break
			}
		}
		fmt.Printf("Finished serving client %s\r\n", con.RemoteAddr().String())
	} else {

	}
}

func (pm *PeerManger) Start() {
	fmt.Println("Getting connections")
	l, err := net.Listen("tcp", pm.ManConfig.Address)
	if err != nil {
		fmt.Println("Failed listening for connections")
		return
	} else {
		for i := 0; i < pm.ManConfig.MaxConnections; i++ {
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
	mc := config.ManagerConfig{
		Address:        "127.0.0.1:6888",
		MaxConnections: 3,
	}
	pc := config.PeerConfig{
		InfoHash:    InfoHash,
		PeerID:      PeerID,
		PieceCount:  128,
		IdleTimeout: 10,
	}
	pm := PeerManger{
		ManConfig:  mc,
		PeerConfig: pc,
	}
	pm.Start()
}
