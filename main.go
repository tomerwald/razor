package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/sha256"
	"log"
	"net"

	"./config"
	"./razor"
	"./tracker"
)

type PeerManger struct {
	ManConfig  config.ManagerConfig
	PeerConfig config.PeerConfig
	Cipher     cipher.Block
}

func (pm *PeerManger) handleClient(con net.Conn) {
	log.Printf("Handling client %s\r\n", con.RemoteAddr().String())
	rc := razor.NewRazorClient(con, &pm.PeerConfig, &pm.Cipher)
	defer rc.Peer.Disconnect()
	if rc.Peer.PerformHandshake() {
		rc.Serve()
	}
}

func (pm *PeerManger) Start() {
	log.Println("Getting connections")
	l, err := net.Listen("tcp", pm.ManConfig.Address)
	if err != nil {
		log.Println("Failed listening for connections")
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
func (pm *PeerManger) StartBC() {
	for _, trackerIP := range pm.ManConfig.Trackers {
		log.Printf("Quering tracker: %s\n", trackerIP)
		con, ResError := net.Dial("udp", trackerIP)
		tc := tracker.TrackerClient{con, 9, 0}
		if ResError == nil {
			err := tc.Connect()
			if err == nil {
				log.Printf("Sending announce: %s\n", trackerIP)
				tc.Announce(pm.PeerConfig.InfoHash, pm.PeerConfig.PeerID)
			}
		}
		address := "127.0.0.1:6882"
		con, serr := net.Dial("tcp", address)
		if serr != nil {
			log.Printf("error connecting to peer: %s", address)
		} else {
			pm.handleClient(con)
		}
	}
}
func main() {
	InfoHash := sha1.Sum([]byte("test"))
	PeerID := []byte("-UW109K-LMYpj9A)8X0R")
	key := sha256.Sum256([]byte("key"))
	enc, _ := aes.NewCipher(key[:])
	mc := config.ManagerConfig{
		Address:        "127.0.0.1:6888",
		MaxConnections: 100,
		Trackers:       []string{"tracker.tiny-vps.com:6969"},
	}
	pc := config.PeerConfig{
		InfoHash:    InfoHash[:],
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
	pm.StartBC()
}
