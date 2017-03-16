package main

import (
	"mydev/pool/pool"
	"fmt"
	"sync"
	"time"
	"github.com/cihub/seelog"
	"os"
)

const (
	CLIENT_CNT = 5000
	LOGGER_CONF_PATH = "D:\\www\\golanxin\\src\\mydev\\pool\\conf\\logger.xml"
)

// TODO
// limit
// 使用map做消息队列？
// wg
// pending 操作的时机

func init() {
	logger, err := seelog.LoggerFromConfigAsFile(LOGGER_CONF_PATH)
	if err != nil {
		seelog.Critical("err parsing config log file ", err)
		os.Exit(1)
	}

	seelog.ReplaceLogger(logger)
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
				OpCode:0,
				RespChan: respCh,
				Data:[]byte(fmt.Sprintf("%d", idx)),
			}
			p.PushJob(fmt.Sprintf("job %d", idx), job)
			res := <-respCh // 等结果
			seelog.Debugf("client %s got result %s\n", idx, res.GetMsg())
		}(i)
	}

	time.Sleep(5*time.Second)
}
