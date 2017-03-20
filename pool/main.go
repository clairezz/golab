package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/cihub/seelog"
	"golab/util"
	"golab/pool/pool"
)

const (
	CLIENT_CNT       = 50
)

// TODO
// limit
// 使用map做消息队列？
// wg
// pending 操作的时机

func init() {
	util.InitLogger()
}

func main() {
	wg := &sync.WaitGroup{}
	p := new(pool.Pool)
	p.Init(wg)
	for i := 0; i < CLIENT_CNT; i++ {
		wg.Add(1)
		go func(idx int) {
			// 生产者线程
			defer wg.Done()
			respCh := make(chan pool.ResultType, 1)
			job := &pool.Job{
				OpCode:   0,
				RespChan: respCh,
				Data:     []byte(fmt.Sprintf("%d", idx)),
			}
			p.PushJob(fmt.Sprintf("job %d", idx), job)
			res := <-respCh // TODO 等结果
			seelog.Debugf("client %s got result %s\n", idx, res.GetMsg())
		}(i)
	}

	time.Sleep(5 * time.Second)
}
