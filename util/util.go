package util

import (
	"sync"
)

// ResultType
type ResultType struct {
	code int
	msg  string
	date []byte
}

func (r *ResultType) GetMsg() string {
	return r.msg
}

// SafePending
type SafePending struct {
	Pending int
	sync.RWMutex
}

func (s *SafePending) Inc() {
	s.Lock()
	s.Pending++
	s.Unlock()
}

func (s *SafePending) Dec() {
	s.Lock()
	s.Pending--
	s.Unlock()
}

func (s *SafePending) Get() int {
	s.RLock()
	n := s.Pending
	s.RUnlock()
	return n
}
