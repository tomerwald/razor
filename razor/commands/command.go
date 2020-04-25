package commands

import (
	"encoding/binary"
)

const (
	Exec   = 1
	Upload = 2
)

type commandError struct {
	message string
}

func (c commandError) Error() string {
	return c.message
}

func NewCommandError(msg string) error {
	return &commandError{message: msg}
}

type Command struct {
	Type    uint32
	Payload []byte
}

func (c *Command) Buffer() []byte {
	metadata := make([]byte, 8)
	binary.BigEndian.PutUint32(metadata[0:4], c.Type)
	binary.BigEndian.PutUint32(metadata[4:8], uint32(len(c.Payload)))
	return append(metadata, c.Payload...)

}

func (c *Command) BufferLen() int {
	return 8 + len(c.Payload)
}

func ReadCommand(buf []byte) ([]byte, error) {
	commandTypeField := buf[0:4]
	LengthField := buf[4:8]
	commandType := binary.BigEndian.Uint32(commandTypeField)
	messageLength := binary.BigEndian.Uint32(LengthField)
	payload := buf[8 : messageLength+8]
	switch commandType {
	case Exec:
		return RunExec(payload)
	case Upload:
		return SaveFile(payload)
	default:
		return nil, NewCommandError("Unknown command type: " + string(commandType))
	}
}
