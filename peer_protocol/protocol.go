package peer_protocol

const (
	Protocol = "\x13BitTorrent protocol"
)

const (
	// BEP 3
	KeepAlive     byte = 99
	Choke         byte = 0
	Unchoke       byte = 1
	Interested    byte = 2
	NotInterested byte = 3
	Have          byte = 4
	Bitfield      byte = 5
	Request       byte = 6
	Piece         byte = 7
	Cancel        byte = 8
	Port          byte = 9

	// BEP 6 - Fast extension
	Suggest     byte = 0x0d // 13
	HaveAll     byte = 0x0e // 14
	HaveNone    byte = 0x0f // 15
	Reject      byte = 0x10 // 16
	AllowedFast byte = 0x11 // 17

	// BEP 10
	Extended byte = 20
)

const (
	HandshakeExtendedID = 0
	RequestMetadataExtensionMsgType = 0
	DataMetadataExtensionMsgType    = 1
	RejectMetadataExtensionMsgType  = 2
)
