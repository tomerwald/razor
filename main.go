package main

import (
	"./peer_protocol"
	"fmt"
	"net"
)

func handleClient(con net.Conn) {
	InfoHash := []byte{0xd5, 0x5f, 0x1e, 0x84, 0x0f, 0x1b, 0xd6,
		0x57, 0x6e, 0xad, 0x67, 0xa4, 0xd0, 0x4e, 0x5d, 0x6e, 0xa2, 0x94, 0x41, 0x4b}
	PeerID := []byte("-UW109K-LMYpj9A)8X0R")
	fmt.Println("Handling client")
	p := peer_protocol.PeerConnection{
		PeerID:   PeerID,
		InfoHash: InfoHash,
		Conn:     con,
	}
	p.Handshake()
}

func connection_loop() {
	fmt.Println("Getting connections")
	l, err := net.Listen("tcp", "127.0.0.1:6888")
	if err != nil {
		fmt.Println("Failed listening for connections")
		return
	} else {
		for i := 0; i < 3; i++ {
			conn, err := l.Accept()
			if err != nil {
				continue
			}
			go handleClient(conn)
		}
	}
}
func main() {
	connection_loop()
}
