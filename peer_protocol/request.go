package peer_protocol

import "encoding/binary"

type RequestMessage struct {
	PieceIndex  uint32
	BlockOffset uint32
	BlockSize   uint32
}

func (r *RequestMessage) Message() Message {
	indexField := make([]byte, 12)
	binary.BigEndian.PutUint32(indexField[0:4], r.PieceIndex)
	binary.BigEndian.PutUint32(indexField[4:8], r.BlockOffset)
	binary.BigEndian.PutUint32(indexField[8:12], r.BlockSize)
	return Message{
		Type:    Request,
		Payload: indexField,
	}
}
func (r RequestMessage) BlockEnd() uint32 {
	return r.BlockOffset + r.BlockSize
}

func NewRequest(PieceIndex uint32, fromOffset uint32, BlockSize uint32) RequestMessage {
	return RequestMessage{
		PieceIndex:  PieceIndex,
		BlockOffset: fromOffset,
		BlockSize:   BlockSize,
	}
}
func RequestFromPayload(payload []byte) RequestMessage {
	return RequestMessage{
		PieceIndex:  binary.BigEndian.Uint32(payload[0:4]),
		BlockOffset: binary.BigEndian.Uint32(payload[4:8]),
		BlockSize:   binary.BigEndian.Uint32(payload[8:12]),
	}
}
