package discovery

import "encoding/json"

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
