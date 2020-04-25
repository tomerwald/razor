package main

import (
	"./config"
	"./razor"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"net"
)

type PeerManger struct {
	ManConfig  config.ManagerConfig
	PeerConfig config.PeerConfig
	Cipher     cipher.Block
}

func (pm *PeerManger) handleClient(con net.Conn) {
	fmt.Printf("Handling client %s\r\n", con.RemoteAddr().String())
	rc := razor.NewRazorClient(con, &pm.PeerConfig, &pm.Cipher)
	defer rc.Peer.Disconnect()
	if rc.Peer.PerformHandshake() {
		rc.Serve()
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
	key := sha256.Sum256([]byte("key"))
	enc, _ := aes.NewCipher(key[:])
	mc := config.ManagerConfig{
		Address:        "127.0.0.1:6888",
		MaxConnections: 100,
	}
	pc := config.PeerConfig{
		InfoHash:    InfoHash,
		PeerID:      PeerID,
		PieceCount:  128,
		IdleTimeout: 120,
		BlockSize:   16384,
		PieceSize:   51 * 16384,
	}
	pm := PeerManger{
		ManConfig:  mc,
		PeerConfig: pc,
		Cipher:     enc,
	}
	pm.Start()
}
