package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/sha256"
	"log"
	"net"
	"time"

	"./config"
	"./razor"
	"./tracker"
)

// PeerManager will manage the tracking and connection to controller peers
type PeerManager struct {
	ManConfig  config.ManagerConfig
	PeerConfig config.PeerConfig
	Cipher     cipher.Block
}

func (pm *PeerManager) handleClient(con net.Conn) {
	log.Printf("Handling client %s\r\n", con.RemoteAddr().String())
	rc := razor.NewRazorClient(con, &pm.PeerConfig, &pm.Cipher)
	defer rc.Peer.Disconnect()
	if rc.Peer.PerformHandshake() {
		rc.Serve()
	}
}

func (pm *PeerManager) connectToControllers(controllers []string) {
	for _, ip := range controllers {
		con, err := net.Dial("tcp", ip)
		if err != nil {
			log.Printf("error connecting to peer: %s", ip)
		} else {
			go pm.handleClient(con)
		}
	}
}
func (pm *PeerManager) getControllers() []string {
	var controllers []string
	for _, trackerIP := range pm.ManConfig.Trackers {
		log.Printf("Quering tracker: %s\n", trackerIP)
		con, ResError := net.Dial("udp", trackerIP)
		tc := tracker.Client{con, 9, 0}
		if ResError == nil {
			err := tc.Connect()
			if err == nil {
				log.Printf("Sending announce: %s\n", trackerIP)
				announceResult, AnounceErr := tc.Announce(pm.PeerConfig.InfoHash, pm.PeerConfig.PeerID)
				if AnounceErr != nil {
					log.Println("Announce request failed")
				} else {
					controllers = append(controllers, announceResult.GetControllerPeers()...)
				}
			}
		}
	}
	return controllers
}

// StartBC will start periodicly searching for controller peers
func (pm *PeerManager) StartBC() {
	for i := 0; i < 100; i++ {
		controllers := pm.getControllers()
		log.Printf("Found %d controllers", len(controllers))
		pm.connectToControllers(controllers)
		time.Sleep(pm.ManConfig.AnnounceInterval)
	}
}

func main() {
	InfoHash := sha1.Sum([]byte("test"))
	PeerID := []byte("-UW109K-LMYpj9A)8X0R")
	key := sha256.Sum256([]byte("key"))
	enc, _ := aes.NewCipher(key[:])
	mc := config.ManagerConfig{
		Address:          "127.0.0.1:6888",
		AnnounceInterval: 5000 * time.Millisecond,
		Trackers:         []string{"127.0.0.1:6969"},
	}
	pc := config.PeerConfig{
		InfoHash:    InfoHash[:],
		PeerID:      PeerID,
		PieceCount:  128,
		IdleTimeout: 120,
		BlockSize:   16384,
		PieceSize:   51 * 16384,
	}
	pm := PeerManager{
		ManConfig:  mc,
		PeerConfig: pc,
		Cipher:     enc,
	}
	pm.StartBC()
}
