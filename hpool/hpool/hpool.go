package hpool

import (
	"container/heap"
	"github.com/juju/errors"

	"golab/hpool/hworker"
	"golab/hpool/hjob"

)

// 管理worker， 每次返回给用户一个负载最小的worker

const (
	WORKER_CNT = 5
)

type Hpool struct {
	workers []*hworker.Hworker
}

func NewHpool () *Hpool{
	p := &Hpool{
		workers: make([]*hworker.Hworker, 0,  WORKER_CNT), //
	}
	for i := 0; i < WORKER_CNT; i++ {
		worker := hworker.NewHworker()
		worker.Run()
		p.workers = append(p.workers, worker)
	}

	return p
}

// Len is the number of elements in the collection.
func (p Hpool) Len() int {
	return len(p.workers)
}
// Less reports whether the element with
// index i should sort before the element with index j.
func (p Hpool) Less(i, j int) bool { // 根据pending判断worker的负载，那么add channel或del channel中的呢？
	return p.workers[i].Pending.Get() < p.workers[j].Pending.Get()
}
// Swap swaps the elements with indexes i and j.
func (p Hpool) Swap(i, j int) {
	p.workers[i], p.workers[j] = p.workers[j], p.workers[i]
}

// add x as element Len()
func (p *Hpool) Push(item interface{}) {
	worker := item.(*hworker.Hworker)
	p.workers = append(p.workers, worker)
}

// remove and return element Len() - 1.
func (p *Hpool) Pop() interface{}  {
	n := p.Len()
	if p.Len() <= 0 {
		return nil
	}
	old := p.workers
	item := old[n -1]
	p.workers = old[:n -1]
	return item
}

func (p *Hpool) PushJob(job *hjob.Hjob) (string, error) {
	if len(p.workers) == 0 {
		return "", errors.New("error: no worker in the pool")
	}

	worker := p.workers[0]
	key := worker.AddJob(job)
	heap.Fix(p,0)

	return key, nil
}