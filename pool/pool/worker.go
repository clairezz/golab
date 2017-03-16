package pool

import (
	"github.com/cihub/seelog"
	"sync"
	"time"
)

type Worker struct {
	index     int // 该参数由调用者在将其放入Pool时指定
	jobqueue  map[string]*Job
	pending   SafePending // 有多少待处理job，
	broadcast chan *Data
	addjob    chan *JobPair
	deljob    chan string
	done      chan struct{}
}

// TODO
type Data struct {
}

func NewWorker(index int, jobqueue_limit int, req_limit int) *Worker {
	return &Worker{
		index:    index,
		jobqueue: make(map[string]*Job, jobqueue_limit),
		pending:  SafePending{0, sync.RWMutex{}},
		//broadcast: make(chan *Data), // TODO
		addjob: make(chan *JobPair, req_limit),
		deljob: make(chan string, req_limit),
		done:   make(chan struct{}),
	}
}

func (w *Worker) Stop() {
	go func() {
		w.done <- struct{}{}
	}()
}

func (w *Worker) AddJob(name string, job *Job) {
	w.addjob <- &JobPair{
		name: name,
		job:  job,
	}
	w.pending.Inc() // !!! 在这里就执行+1操作，而不是等到真正添加到map之后,否则负载均衡效果不好
}

func (w *Worker) RemoveJob(name string) {
	w.deljob <- name
	w.pending.Dec()
	// TODO 在哪里进行-1操作，如果移除不成功怎么办
}

func (w *Worker) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	// 以下goroutine是为了实现线程安全的map
	wg.Add(1)
	go func() {
		seelog.Debugf("New goroutine, worker index: %d\n", w.index)
		defer wg.Done()
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for { // 实现线程安全的map
			select {
			case jobPaire := <-w.addjob:
				w.insertJob(jobPaire.name, jobPaire.job)
			case jobname := <-w.deljob:
				w.deleteJob(jobname)
			case <-ticker.C:
			//	seelog.Debugf("worker %d !\n", w.index) // 定期报告自己
			case <-w.done:
				seelog.Debugf("worker %d exit!\n", w.index)
				break
			}
		}
	}()

	// worker开始工作
	for {
		for key, job := range w.jobqueue {
			w.RemoveJob(key)
			seelog.Debugf("worker %d is doing job: %s, whose pending is %d\n", w.index, string(job.Data), w.pending.Get())
			code2op[job.OpCode](job.Data, job.RespChan)
			seelog.Debugf("worker %d done job: %s\n, whose pending is %d\n", w.index, string(job.Data), w.pending.Get())
		}
	}

}

// TODO 重复key
// 实际的map操作：添加
func (w *Worker) insertJob(key string, job *Job) {
	w.jobqueue[key] = job
}

// TODO 要删除的job还在addChan中或者根本不存在
// 实际的map操作：删除
func (w *Worker) deleteJob(jobname string) {
	delete(w.jobqueue, jobname)
}

func (w *Worker) GetPending() int {
	return w.pending.Get()
}

func (w *Worker) GetIndex() int {
	return w.index
}
