package serviceManage

import "go-im/ipConfig/source"

type Candidate struct {
	IP    string
	Port  string
	Score float64
}

func NewCandidate(e *source.Event, score float64) *Candidate {
	return &Candidate{
		IP:    e.Ip,
		Port:  e.Port,
		Score: score,
	}
}
