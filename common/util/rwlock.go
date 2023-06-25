package util

import "sync"

type RWMutex struct {
	lock   sync.Mutex
	wLock  sync.Mutex
	reader int
	writer int
}

// RLock 获取读锁
func (m *RWMutex) RLock() {
	m.lock.Lock()
	m.reader++
	if m.reader == 1 {
		//第一个获取读锁，需要获取写锁
		m.wLock.Lock()
	}
	m.lock.Unlock()
}

// RUnLock 释放读锁
func (m *RWMutex) RUnLock() {
	m.lock.Lock()
	m.reader--
	if m.reader == 0 {
		m.wLock.Unlock()
	}
	m.lock.Unlock()
}

func (m *RWMutex) Lock() {
	m.lock.Lock()
	m.writer++
	m.lock.Unlock()
	m.wLock.Lock()
}

func (m *RWMutex) UnLock() {
	m.wLock.Unlock()
}
