package msgStorage

import (
	"github.com/panjf2000/ants/v2"
	"go-im/common/log"
	"go-im/common/message"
)

// msgStorage 消息存储
type msgStorage struct {
	onlineMsg  chan *message.Cmd // 在线消息
	offLineMsg chan *message.Cmd // 离线消息
	done       chan struct{}
	workPool   *ants.Pool //协程池
}

// initMsgStorage  初始化消息存储器
func initMsgStorage() (*msgStorage, error) {
	pool, err := ants.NewPool(1000)
	if err != nil {
		return nil, err
	}
	return &msgStorage{
		onlineMsg:  make(chan *message.Cmd, 10),
		offLineMsg: make(chan *message.Cmd, 10),
		done:       make(chan struct{}),
		workPool:   pool,
	}, nil
}

// msgProc 消息处理
func (s *msgStorage) msgProc() {
	select {
	case onMsg := <-s.onlineMsg:
		s.handleOnlineMsg(onMsg)
	}
}

// 处理现在消息
func (s *msgStorage) handleOnlineMsg(onMsg *message.Cmd) {
	for {
		if err := s.workPool.Submit(
			func() {
				//TODO 消息转发给用户

				//TODO 消息持久化
			},
		); err != nil {
			log.Info(err)
			continue
		}
		break
	}
}
