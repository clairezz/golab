package hworker

import (
	"golab/hpool/hjob"
	"golab/pool/pool"
	"github.com/satori/go.uuid"
	"sync"
	"github.com/cihub/seelog"
)

const (
	ADDCHAN_BUF = 2
	DELCHAN_BUF = 2
	JOBQUEUE_LEN = 10 // TODO 如何限制？ ==》由 pool 控制？
)

// hworker 接收 job 并执行job, 需要实现并发安全
// 并发安全的 map : 增删给查都交给同一个 goroutine 处理
type HjobWrap struct {
	key   string
	value *hjob.Hjob
}

type Hworker struct {
	jq      map[string] *hjob.Hjob // 任务队列， 以任务名为 key
	Pending pool.SafePending       // 待执行的任务数量， 包括 channel 中的（还没来得及 insert 进 map）和 map 中的
	addChan chan HjobWrap          // 向 map 中添加一个元素， sender 是用户， receiver 是执行任务的那个 goroutine
	delChan chan string            // 从 map 中删除一个元素， sender 是用户
}

func NewHworker () *Hworker{
	return &Hworker{
		jq: make(map[string] *hjob.Hjob, JOBQUEUE_LEN),
		Pending: pool.SafePending{0,sync.RWMutex{}},
		addChan: make(chan HjobWrap, ADDCHAN_BUF),
		delChan: make(chan string, DELCHAN_BUF),
	}
}

// 为避免重名，用户不用提供jobName，由函数实现生成uuid作为jobName，并返回jobName， 以便于删除等操作
func (hw *Hworker) AddJob (job *hjob.Hjob) string {
	key := uuid.NewV4().String()
	seelog.Debugf("uuid key: %s\n", key)
	jw := HjobWrap{key, job}
	hw.addChan <- jw
	hw.Pending.Inc()

	return key
}

func (hw *Hworker) DelJob (jobName string) {
	hw.delChan <- jobName
	hw.Pending.Dec() // TODO 一定能删除成功吗？ 实际没有删除成功pending却减 1 的情况如何处理？
}

func (hw *Hworker) Run () {
	go func() {
		for {
			select {
			case jw := <-hw.addChan: // 添加信号
				hw.jq[jw.key] = jw.value // 实际的添加
			case jn := <-hw.delChan: // 删除信号
				hw.addAll() // 先将addChan中的所有job添加到map中再执行delete，避免delete时还未add
				delete(hw.jq, jn) // 实际的删除
			default: // 取出一个job并执行
				for jobName, job := range hw.jq { // 查询一个元素
					delete(hw.jq, jobName)
					job.Handler(job.Data, job.RespCh)

					break
				}
			}
		}
	}()
}

func (hw *Hworker) addAll() {
	for i := 0; i < len(hw.addChan); i++ {
		jw := <- hw.addChan
		hw.jq[jw.key] = jw.value
	}
}