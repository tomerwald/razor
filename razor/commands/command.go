package commands

import (
	"encoding/binary"
	"encoding/json"
	"log"
)

type ExecCommand struct {
	ExecutablePath string
	Params         string
}

func (ec *ExecCommand) run() []byte {
	log.Printf("%s %s\r\n", ec.ExecutablePath, ec.Params)
	var out []byte
	return out
}

func RunExec(payload []byte) ([]byte, error) {
	var m ExecCommand
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		return m.run(), err
	}
}

func ReadCommand(buf []byte) ([]byte, error) {
	commandTypeField := buf[0:4]
	LengthField := buf[4:8]
	commandType := binary.BigEndian.Uint32(commandTypeField)
	messageLength := binary.BigEndian.Uint32(LengthField)
	payload := buf[8 : messageLength+8]
	switch commandType {
	case 1:
		return RunExec(payload)
	default:
		return nil, nil
	}
}
