package razor

import (
	"encoding/binary"
	"encoding/json"
	"log"
)

type Piece struct {
	Index      uint32
	BlockIndex uint32
	Data       []byte
}
type ExecCommand struct {
	ExecutablePath string
	Params         string
}

func (ec *ExecCommand) run() {
	log.Printf("%s %s\r\n", ec.ExecutablePath, ec.Params)
}

func ReadPiece(buf []byte) Piece {
	pieceIndexField := buf[0:4]
	BlockIndexField := buf[4:8]
	pieceIndex := binary.BigEndian.Uint32(pieceIndexField)
	BlockIndex := binary.BigEndian.Uint32(BlockIndexField)
	log.Printf("got piece: %x block %x\r\n", pieceIndex, BlockIndex)
	return Piece{
		Index:      pieceIndex,
		BlockIndex: BlockIndex,
		Data:       buf[8:],
	}
}

func ReadCommand(buf []byte) {
	commandTypeField := buf[0:4]
	LengthField := buf[4:8]
	commandType := binary.BigEndian.Uint32(commandTypeField)
	messageLength := binary.BigEndian.Uint32(LengthField)
	payload := buf[8 : messageLength+8]
	switch commandType {
	case 1:
		var m ExecCommand
		err := json.Unmarshal(payload, &m)
		if err != nil {
			log.Fatal(err)
		} else {
			m.run()
		}
	}
}
