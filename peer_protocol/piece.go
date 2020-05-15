package peer_protocol

import (
	"encoding/binary"
	"log"
)

type PieceMessage struct {
	// a message containing data for a requested piece
	PieceIndex  uint32
	BlockOffset uint32
	Data        []byte
}

func (p *PieceMessage) Message() Message {
	// create a message of Piece type
	metadataField := make([]byte, 8)
	binary.BigEndian.PutUint32(metadataField[0:4], p.PieceIndex)
	binary.BigEndian.PutUint32(metadataField[4:8], p.BlockOffset)
	metadataField = append(metadataField, p.Data...)
	return Message{
		Type:    Piece,
		Payload: metadataField,
	}
}

func ReadPiece(buf []byte) PieceMessage {
	pieceIndexField := buf[0:4]
	BlockIndexField := buf[4:8]
	pieceIndex := binary.BigEndian.Uint32(pieceIndexField)
	BlockIndex := binary.BigEndian.Uint32(BlockIndexField)
	log.Printf("got piece: %x block %x\r\n", pieceIndex, BlockIndex)
	return PieceMessage{
		PieceIndex:  pieceIndex,
		BlockOffset: BlockIndex,
		Data:        buf[8:],
	}
}
