package serviceManage

import (
	"github.com/bytedance/gopkg/util/logger"
	"go-im/ipConfig/source"
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
func DisPatch() []Candidate {
	c := d.getCandidates()
	if len(c) <= 5 {
		return c
	}
	return c[:5]
}

func (d *DisPatcher) getCandidates() []Candidate {
	d.rwMu.RLock()
	defer d.rwMu.RUnlock()
	c := make([]Candidate, len(d.candidateTable))
	for _, v := range d.candidateTable {
		c = append(c, *v)
	}
	return c
}

func (d *DisPatcher) addCandidate(event *source.Event) {
	d.rwMu.Lock()
	defer d.rwMu.Unlock()
	d.candidateTable[event.Key()] = NewCandidate(event)
	logger.Info("add candidate ", event.Key())
}

func (d *DisPatcher) delCandidate(event *source.Event) {
	d.rwMu.Lock()
	defer d.rwMu.Unlock()
	delete(d.candidateTable, event.Key())
	logger.Info("del candidate ", event.Key())
}
