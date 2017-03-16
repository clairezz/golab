package pool

import (
	"container/heap"
	"sync"
	"github.com/cihub/seelog"
	"time"
)

const (
	WORKER_CNT = 5
	JOBQUEUE_LIMIT = 10
	REQ_LIMIT = 10
)

// 实现以下5个接口
type Pool struct {
	workers []*Worker;
	sync.Mutex
}

// Len is the number of elements in the collection.
func (p Pool) Len() int {
	return len(p.workers)
}
// Less reports whether the element with
// index i should sort before the element with index j.
func (p Pool) Less(i, j int) bool { // 根据pending判断worker的负载，那么add channel或del channel中的呢？
	return p.workers[i].pending.Get() < p.workers[j].pending.Get()
}
// Swap swaps the elements with indexes i and j.
func (p Pool) Swap(i, j int) {
	p.workers[i], p.workers[j] = p.workers[j], p.workers[i]
	p.workers[i].index = i
	p.workers[j].index = j
}

// add x as element Len()
func (p *Pool) Push(x interface{}) {
	n := p.Len()
	item := x.(*Worker)
	p.workers = append(p.workers, item)
	item.index = n
}

// remove and return element Len() - 1.
func (p *Pool) Pop() interface{}  {
	n := p.Len()
	if p.Len() <= 0 {
		return nil
	}
	old := p.workers
	item := old[n -1]
	p.workers = old[:n -1]
	return item
}

func (p *Pool) WorkerCnt() int {
	return len(p.workers)
}
func (p *Pool) Init(wg *sync.WaitGroup) {
	for i := 0; i < WORKER_CNT; i++ {
		worker := NewWorker(i,JOBQUEUE_LIMIT,REQ_LIMIT)
		go worker.Run(wg)
		seelog.Debugf("heap.Push worker Index: %d, Pending: %d\n", worker.GetIndex(), worker.GetPending())
		heap.Push(p, worker)
	}
	time.Sleep(1*time.Second)
	heap.Init(p)
}

func (p *Pool) PushJob(jobname string, job *Job) {
	p.Lock() // !!!
	defer p.Unlock()
	worker := p.workers[0]
	seelog.Debugf("Add job %s to worker %d, whose pending is %d\n", string(job.Data), worker.index, worker.pending.Get())
	worker.AddJob(jobname,job)
	heap.Fix(p, worker.index)
	//heap.Init(p)
	seelog.Debugf("%s:---------------------------pending range-----------------------\n", string(job.Data))
	for _, w := range p.workers {
		seelog.Debugf("%s: worker %d has pending %d\n", string(job.Data), w.index, w.pending.Get())
	}
	seelog.Debugf("%s:---------------------------pending range end--------------------\n", string(job.Data))
}