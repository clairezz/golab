package pool

import (
	"sync"
	"fmt"
)

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

func Hello(data []byte, resCh chan ResultType) {
	resCh <- ResultType{
		code: 0,
		msg:  "Hello secceed" + string(data[:]),
	}
	res := <- resCh
	fmt.Println(res.GetMsg())
}

