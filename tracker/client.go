package tracker

import (
	"net"

	"./requests"
)

type Client struct {
	Con           net.Conn
	TransactionID uint32
	ConnectionID  uint64
}

func (t *Client) readConnectResponse() (uint64, error) {
	buf := make([]byte, 16)
	t.Con.Read(buf)
	conResponse, err := requests.UnmarshalConnectResponse(buf)
	return conResponse.ConnectionID, err
}

func (t *Client) Connect() error {
	conRequest := requests.ConnectionRequest()
	requestBuffer := conRequest.MarshalBinary()
	if _, err := t.Con.Write(requestBuffer); err != nil {
		return err
	}
	connID, err := t.readConnectResponse()
	if err != nil {
		return err
	}
	t.ConnectionID = connID
	return nil
}

func (t *Client) Announce(infoHash []byte, peerID []byte) (requests.AnnounceResponse, error) {
	requestBuffer := requests.CreateAnnounceRequest(infoHash, peerID, t.ConnectionID)
	t.Con.Write(requestBuffer.MarshalBinary())
	resultBuffer := make([]byte, 320)
	if _, err := t.Con.Read(resultBuffer); err != nil {
		return requests.AnnounceResponse{}, requests.NewAnnouncementError("Could not read announce response")
	}
	res, parseError := requests.UnmarshalAnnounceResponse(resultBuffer)
	if parseError != nil {
		return requests.AnnounceResponse{}, parseError
	}
	return res, nil
}
