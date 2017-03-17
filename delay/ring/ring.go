package ringr

import (
	"container/ring"
	"time"

	"github.com/cihub/seelog"
)

// TODO
// task中存放回调函数
// AddTask线程安全：使用线程安全的map
// 单测

const (
	DURATION = 6 // 环指针每DURATION秒钟走一格
	TOTAL_DUR = 36 // 一圈的总时间（以秒为单位）, 为了便于测试 置为36秒
)

type Loop struct{
	r *ring.Ring
	done chan struct{}
}

type Task struct {
	t int // 该任务t秒之后被执行
	msg string
}

func NewTask(t int, msg string) *Task {
	return &Task{t,msg}
}

type taskWrap struct {
	task      Task
	cycle_cnt int
}

type node struct {
	taskSet map[*taskWrap]bool
	index int // 该节点在环上的编号
}

func (n node) Index() int {
	return n.index
}

func NewLoop() *Loop{
	n := TOTAL_DUR / DURATION // 为了便于测试
	l := &Loop{
		r:	ring.New(n),
		done:	make(chan struct{}),
	}

	for i := 0; i < l.r.Len(); i++ {
		l.r.Value = node{
			taskSet: make(map[*taskWrap]bool,0),
			index: i,
		}
		l.r = l.r.Next()
	}
	return l
}

func (l *Loop) Len() int {
	return l.r.Len()
}

func (l *Loop) AddTask(task Task) { // 多协程程操作map不安全，TODO 线程安全的map
	d := task.t / DURATION % l.r.Len() // !!!
	n := task.t /DURATION/ l.r.Len() // !!!
	nd := &taskWrap{
		task:task,
		cycle_cnt:n,
	}

	dst := l.r.Move(d)
	dst.Value.(node).taskSet[nd] = true
}

func (l *Loop) AddTsks(tasks []Task) {
	for _, task := range tasks {
		l.AddTask(task)
	}
}

func (l *Loop) Run() {
	go func() {
		ticker := time.NewTicker(DURATION * time.Second)
		seelog.Debugf("\t^^^^^loop starts! it ticks every ", DURATION, " s ^^^^^")
		defer ticker.Stop()
		for{
			select {
			case <-l.done:
				seelog.Debugf("\t$$$$$\tloop stoped! bye\t$$$$$") // 收到结束信号
				break
			case <-ticker.C:
				l.do() // 时间到，执行当前结点的任务
			}
		}
	}()
}

func (l *Loop) Stop() {
	l.done <- struct {}{}
}

func (l *Loop) Status() {
	seelog.Debugf("report start--------------------report start")
	seelog.Debugf("Loop Len: %d", l.r.Len())
	seelog.Debugf("current node index: %d", l.r.Value.(node).Index())
	rr := l.r // 复制一份当前指针，用该指针遍历环
	for i := 0; i < rr.Len(); i++ {
		for t := range rr.Value.(node).taskSet {
			seelog.Debugf("******node index: %d******task: %s;\t\tcycle_cnt: %d\n", rr.Value.(node).index, t.task.msg, t.cycle_cnt)
		}
		rr = rr.Next() // 移动指针
	}
	seelog.Debugf("report end--------------------report end")
}

func (l *Loop) do() {
	seelog.Debugf("@@@@@@@@@@@@@@@@@@@@@@@\t%s\t@@@@@@@@@@@@@@@@@@@@@@@@", time.Now().Format("2006-01-02 15:04:05.99")) // 当前指针移动一个位置，并检查当前节点是否有需要执行的任务
	seelog.Debugf("--------------------------before execute task--------------------------")
	l.Status()
	for nd := range l.r.Value.(node).taskSet {
		if nd.cycle_cnt == 0 {
			seelog.Debugf("sengding msg: %s", nd.task.msg)// 执行任务
			delete(l.r.Value.(node).taskSet, nd) // 移除任务
		} else {
			nd.cycle_cnt--
		}
	}
	seelog.Debugf("--------------------------after execute task--------------------------")
	l.Status()

	// 方案1：指针指向待处理结点——等时钟响时，执行当前结点任务先执行当前结点的任务，然后再移动指针
	// 方案2：指针指向已处理结点——等时钟响时，先移动指针，然后执行新结点任务
	// 选择方案1，因为对于方案2，初始结点上的任务会被延迟一轮才能执行
	l.r =	l.r.Move(1)
	return
}