package ring

import (
	"container/ring"
	"time"
	"fmt"
)

type Loop struct {
	r *ring.Ring
}

type Task struct {
	t int
	msg string
}

func NewTask(t int, msg string) *Task {
	return &Task{t,msg}
}

type node struct {
	task Task
	n int``
}

/**
 * @param d: 间隔的秒数
 */
func NewLoop(d int) *Loop{
//	n := 3600 / d
	l := &Loop{
		ring.New(d),
	}

	for i := 0; i < l.r.Len(); i++ {
		l.r.Value = make(map[*node]bool,0)
		l.r = l.r.Next()
	}
	return l
}

func (l *Loop) AddTsk(task Task) {
	d := task.t % l.r.Len()
	n := task.t / l.r.Len()
	nd := &node{
		task:task,
		n:n,
	}

	dst := l.r.Move(d)
	dst.Value.(map[*node]bool)[nd] = true
}

func (l *Loop) Run() {
	go func() {
		ticker := time.NewTicker(1*time.Second)
		defer ticker.Stop()
		for{
			<-ticker.C
			fmt.Println("……")
			l.r =	l.r.Move(1)
			for nd := range l.r.Value.(map[*node]bool) {
				if nd.n == 0 {
					fmt.Printf("sengding msg: %s\n", nd.task.msg)// 执行任务
					delete(l.r.Value.(map[*node]bool), nd) // 移除任务
				} else {
					nd.n--
				}
			}
		}
	}()
}
