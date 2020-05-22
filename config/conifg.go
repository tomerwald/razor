package config

import "time"

type ManagerConfig struct {
	Address          string
	AnnounceInterval time.Duration
	Trackers         []string
}
type PeerConfig struct {
	InfoHash    []byte
	PeerID      []byte
	PieceCount  int
	IdleTimeout int
	BlockSize   uint32
	PieceSize   uint32
}
