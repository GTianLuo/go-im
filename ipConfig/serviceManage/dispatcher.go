package serviceManage

import (
	"go-im/ipConfig/source"
	"go-im/log"
	"sort"
	"sync"
)

type DisPatcher struct {
	candidateTable map[string]*Candidate
	rwMu           sync.RWMutex
}

var d *DisPatcher

// Init 初始化监听event事件
func Init() {
	d = &DisPatcher{
		candidateTable: make(map[string]*Candidate),
	}
	go func() {
		eventChan := source.EventChan()
		for event := range eventChan {
			switch event.Type {
			case source.AddNodeEvent:
				d.addCandidate(event)
			case source.DelNodeEvent:
				d.delCandidate(event)
			}
		}
	}()
}

// DisPatch 调度可用服务
func DisPatch() []string {
	c := d.getCandidates()
	sort.Slice(c, func(i, j int) bool {
		return c[i].Score < c[j].Score
	})
	s := make([]string, 0, 5)
	for i := 0; i < 5; i++ {
		if i >= len(c) {
			return s
		}
		s = append(s, c[i].IP+":"+c[i].Port)
	}
	return s
}

func (d *DisPatcher) getCandidates() []Candidate {
	d.rwMu.RLock()
	defer d.rwMu.RUnlock()
	c := make([]Candidate, 0, 1)
	for _, v := range d.candidateTable {
		c = append(c, *v)
	}
	return c
}

func (d *DisPatcher) addCandidate(event *source.Event) {
	d.rwMu.Lock()
	defer d.rwMu.Unlock()
	d.candidateTable[event.Key()] = NewCandidate(event, getScore(event.ConnectNum, event.MessageBytes))
	log.Info("add candidate ", event.Key(), "hold service:", len(d.candidateTable))
}

func (d *DisPatcher) delCandidate(event *source.Event) {
	d.rwMu.Lock()
	defer d.rwMu.Unlock()
	delete(d.candidateTable, event.Key())
	log.Info("del candidate ", event.Key())
}

func getScore(connNums float64, messageBytes float64) float64 {
	return connNums/1024 + messageBytes
}
