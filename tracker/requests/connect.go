package requests

import (
	"encoding/binary"
	"math/rand"
)

const (
	ConnectionID = 0x41727101980
)

type ConnectionError struct {
	message string
}

func (e ConnectionError) Error() string {
	return e.message
}
func NewConnectionError(msg string) ConnectionError {
	return ConnectionError{msg}
}

type Connect struct {
	TransactionID uint32
	ConnectionID  uint64
}

func (c *Connect) MarshalBinary() []byte {
	ConnectBuffer := make([]byte, 16)
	binary.BigEndian.PutUint64(ConnectBuffer, uint64(ConnectionID))
	binary.BigEndian.PutUint32(ConnectBuffer[8:], 0) // Connect action
	binary.BigEndian.PutUint32(ConnectBuffer[12:], c.TransactionID)
	return ConnectBuffer
}
func UnmarshalConnectResponse(buf []byte) (Connect, error) {
	if len(buf) != 16 {
		return Connect{}, NewConnectionError("Wrong sized connect response")
	}
	action := binary.BigEndian.Uint32(buf[:4])
	if action != 0 {
		return Connect{}, NewConnectionError("An error occured")
	}
	transactionID := binary.BigEndian.Uint32(buf[4:8])
	connectionID := binary.BigEndian.Uint64(buf[8:16])
	return Connect{transactionID, connectionID}, nil

}
func ConnectionRequest() Connect {
	return Connect{
		rand.Uint32(),
		0,
	}
}
