package hworker

import (
	"golab/hpool/hjob"
	"golab/pool/pool"
	"github.com/satori/go.uuid"
)

const (
	ADDCHAN_BUF = 1
	DELCHAN_BUF = 1
	JOBQUEUE_LEN = 10 // TODO 如何限制？ ==》由 pool 控制？
)

// hworker 接收 job 并执行job, 需要实现并发安全
// 并发安全的 map : 增删给查都交给同一个 goroutine 处理
type HjobWrap struct {
	key   string
	value *hjob.Hjob
}

type Hworker struct {
	jq map[string] *hjob.Hjob // 任务队列， 以任务名为 key
	pending pool.SafePending // 待执行的任务数量， 包括 channel 中的（还没来得及 insert 进 map）和 map 中的
	addChan chan HjobWrap // 向 map 中添加一个元素， sender 是用户， receiver 是执行任务的那个 goroutine
	delChan chan string // 从 map 中删除一个元素， sender 是用户
}

func NewHworker () {
	return &Hworker{
		jq: make(map[string] *hjob.Hjob, JOBQUEUE_LEN),
		pending: pool.SafePending{0},
		addChan: make(chan HjobWrap, ADDCHAN_BUF),
		delChan: make(chan HjobWrap, DELCHAN_BUF),
	}
}

// 用户不用提供jobName，由函数实现生成uuid作为jobName，并返回jobName， 以便于删除等操作
func (hw *Hworker) AddJob (job *Hworker) string { // TODO 重名job
	key := uuid.NewV4().String()
	hw.addChan <- HjobWrap{key, job}
	hw.pending.Inc()
}

func (hw *Hworker) DelJob (jobName string) {
	hw.delChan <- jobName
	hw.pending.Dec() // TODO 一定能删除成功吗？ 实际没有删除成功pending却减 1 的情况如何处理？
}

func (hw *Hworker) Run () {
	go func() {
		for {
			select {
			case jw := <-hw.addChan: // 添加
				hw.jq[jw.key] = jw.value
			case jn := <-hw.delChan: // 删除
				delete(hw.jq, jn)
			default:
				for jobName, job := range hw.jq { // 查询一个元素
					delete(hw.jq, jobName)
					job.Handler()
					break
				}
			}
		}
	}()
}
