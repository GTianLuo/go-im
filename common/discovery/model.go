package discovery

import (
	"encoding/json"
	"strings"
)

type EndpointInfo struct {
	IP       string
	Port     string
	Metadata map[string]interface{}
}

func (ed *EndpointInfo) Marshal() (string, error) {
	data, err := json.Marshal(ed)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (ed *EndpointInfo) UnMarshal(data []byte) error {
	return json.Unmarshal(data, ed)
}

func Transform(addr string, connNum float64, messageBytes float64) *EndpointInfo {
	s := strings.Split(addr, ":")
	metadate := map[string]interface{}{"connect_num": connNum, "message_bytes": messageBytes}
	return &EndpointInfo{
		IP:       s[0],
		Port:     s[1],
		Metadata: metadate,
	}
}
