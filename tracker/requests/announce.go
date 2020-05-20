package requests

import (
	"encoding/binary"
	"math/rand"
	"net"
	"fmt"
)

type AnnouncementError struct {
	message string
}

func (e AnnouncementError) Error() string {
	return e.message
}
func NewAnnouncementError(msg string) AnnouncementError {
	return AnnouncementError{msg}
}

// AnnounceRequest stuff
type AnnounceRequest struct {
	ConnectionID  uint64
	TransactionID uint32
	InfoHash      []byte
	PeerID        []byte
}

func (ar *AnnounceRequest) MarshalBinary() []byte {
	commonBuffer := make([]byte, 16)
	binary.BigEndian.PutUint64(commonBuffer, uint64(ar.ConnectionID))
	binary.BigEndian.PutUint32(commonBuffer[8:], 1) // Connect action
	binary.BigEndian.PutUint32(commonBuffer[12:], ar.TransactionID)
	commonBuffer = append(commonBuffer, ar.InfoHash...)
	commonBuffer = append(commonBuffer, ar.PeerID...)
	paramBuffer := make([]byte, 42)
	binary.BigEndian.PutUint64(paramBuffer[0:], 1024)  //downloaded
	binary.BigEndian.PutUint64(paramBuffer[8:], 1024)  // left
	binary.BigEndian.PutUint64(paramBuffer[16:], 0)     // uploaded
	binary.BigEndian.PutUint32(paramBuffer[24:], 0)     // event
	binary.BigEndian.PutUint32(paramBuffer[28:], 0)     // IP (commonly unimplemented)
	binary.BigEndian.PutUint32(paramBuffer[32:], 0)     // key
	binary.BigEndian.PutUint32(paramBuffer[36:], 50)    // num want
	binary.BigEndian.PutUint16(paramBuffer[40:], 19238) // port
	return append(commonBuffer, paramBuffer...)
}
func CreateAnnounceRequest(infoHash []byte, peerID []byte, connectionID uint64) AnnounceRequest {
	return AnnounceRequest{
		InfoHash:      infoHash,
		PeerID:        peerID,
		ConnectionID:  connectionID,
		TransactionID: rand.Uint32(),
	}
}

type AnnounceResponse struct {
	Interval uint32
	Leechers uint32
	Seeders  uint32
	Peers    []string
}

func parseIpField(buf []byte) []string {
	addresses := make([]string, len(buf)/6)
	for i := 0; i < len(buf)/6; i++ {
		ip := net.IP(buf[i*6 : (i*6)+4])
		port := binary.BigEndian.Uint16(buf[(i*6)+4 : (i+1)*6])
		if port > 0{
			addresses[i] = fmt.Sprintf("%s:%d", ip.String(), port)
		}
	}
	return addresses
}
func UnmarshalAnnounceResponse(buf []byte) (AnnounceResponse, error) {
	// read a scrape response, contianing the seeders, completed and leechers
	if len(buf) < 20 || len(buf[20:])%6 != 0 {
		return AnnounceResponse{}, NewAnnouncementError("Wrong sized connect response")
	}
	action := binary.BigEndian.Uint32(buf[:4])
	if action != 1 {
		return AnnounceResponse{}, NewAnnouncementError("An error occured")
	}
	Interval := binary.BigEndian.Uint32(buf[8:12])
	leechers := binary.BigEndian.Uint32(buf[12:16])
	seeders := binary.BigEndian.Uint32(buf[16:20])
	addresses := parseIpField(buf[20:])
	return AnnounceResponse{
		Leechers: leechers,
		Interval: Interval,
		Seeders:  seeders,
		Peers:    addresses,
	}, nil
}
