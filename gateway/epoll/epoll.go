package epoll

import (
	"go-im/common/conf/serviceConf"
	"sync"
	"syscall"
)

const EPOLLET = 1 << 31

var epollMap map[int]*Epoller = make(map[int]*Epoller)
var mu sync.Mutex

type Epoller struct {
	Epollfd int
}

func GetEpoller(epollfd int) *Epoller {
	return epollMap[epollfd]
}

func CreateEpoll() (*Epoller, error) {

	//创建非阻塞的epoll
	efd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	epoller := &Epoller{Epollfd: efd}
	mu.Lock()
	epollMap[efd] = epoller
	mu.Unlock()
	return epoller, nil
}

func (e *Epoller) AddEpollTask(fd int32) error {
	events := &syscall.EpollEvent{Events: EPOLLET | syscall.EPOLLIN | syscall.EPOLLHUP, Fd: fd}
	err := syscall.EpollCtl(e.Epollfd, syscall.EPOLL_CTL_ADD, int(fd), events)
	return err
}

func (e *Epoller) DelEpollTask(fd int32) error {
	err := syscall.EpollCtl(e.Epollfd, syscall.EPOLL_CTL_DEL, int(fd), nil)

	return err
}

func (e *Epoller) EventTigger() ([]syscall.EpollEvent, int, error) {
	events := make([]syscall.EpollEvent, serviceConf.GetGateWayEpollMaxTriggerConn())
	// timeout设置为200,避免长时间阻塞等待
	n, err := syscall.EpollWait(e.Epollfd, events, 200)
	if err != nil {
		return events, 0, err
	}
	return events, n, nil
}

func (e *Epoller) Close() error {
	return syscall.Close(e.Epollfd)
}
