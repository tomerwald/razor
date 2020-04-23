package peer_protocol

import (
	"encoding/binary"
	"log"
)

type PieceMessage struct {
	Index      uint32
	BlockIndex uint32
	Data       []byte
}

func ReadPiece(buf []byte) PieceMessage {
	pieceIndexField := buf[0:4]
	BlockIndexField := buf[4:8]
	pieceIndex := binary.BigEndian.Uint32(pieceIndexField)
	BlockIndex := binary.BigEndian.Uint32(BlockIndexField)
	log.Printf("got piece: %x block %x\r\n", pieceIndex, BlockIndex)
	return PieceMessage{
		Index:      pieceIndex,
		BlockIndex: BlockIndex,
		Data:       buf[8:],
	}
}
