package requests

import "encoding/binary"

type Scrape struct {
	InfoHash      []byte
	TransactionID uint32
	ConnectionID  uint64
}

func (s *Scrape) MarshalBinary() []byte {
	scrapeBuffer := make([]byte, 36)
	binary.BigEndian.PutUint64(scrapeBuffer, uint64(s.ConnectionID))
	binary.BigEndian.PutUint32(scrapeBuffer, uint32(2)) // scrape action
	binary.BigEndian.PutUint32(scrapeBuffer, uint32(s.TransactionID))
	scrapeBuffer = append(scrapeBuffer, s.InfoHash...)
	return scrapeBuffer
}

func UnmarshalScrapeResponse(buf []byte) (uint32, uint32, uint32, error) {
	// read a scrape response, contianing the seeders, completed and leechers
	if len(buf) != 20 {
		return 0, 0, 0, NewConnectionError("Wrong sized connect response")
	}
	action := binary.BigEndian.Uint32(buf[:4])
	if action != 2 {
		return 0, 0, 0, NewConnectionError("An error occured")
	}
	seeders := binary.BigEndian.Uint32(buf[8:12])
	completed := binary.BigEndian.Uint32(buf[12:16])
	leechers := binary.BigEndian.Uint32(buf[16:20])
	return seeders, completed, leechers, nil
}
