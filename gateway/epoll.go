package gateway

import (
	"go-im/common/conf/serviceConf"
	"syscall"
)

const EPOLLET = 1 << 31

type epoller struct {
	epollfd int
}

func createEpoll() (*epoller, error) {
	//创建非阻塞的epoll
	efd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoller{epollfd: efd}, nil
}

func (e *epoller) addEpollTask(fd int32) error {
	events := &syscall.EpollEvent{Events: EPOLLET | syscall.EPOLLIN | syscall.EPOLLHUP, Fd: fd}
	err := syscall.EpollCtl(e.epollfd, syscall.EPOLL_CTL_ADD, int(fd), events)
	return err
}

func (e *epoller) delEpollTask(fd int32) error {
	err := syscall.EpollCtl(e.epollfd, syscall.EPOLL_CTL_DEL, int(fd), nil)

	return err
}

func (e *epoller) eventTigger() ([]syscall.EpollEvent, int, error) {
	events := make([]syscall.EpollEvent, serviceConf.GetGateWayEpollMaxTriggerConn())
	n, err := syscall.EpollWait(e.epollfd, events, 200)
	if err != nil {
		return events, 0, err
	}
	return events, n, nil
}

func (e *epoller) close() error {
	return syscall.Close(e.epollfd)
}
