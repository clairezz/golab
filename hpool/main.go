package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/cihub/seelog"

	"github.com/clairezz/golab/hpool/hjob"
	"github.com/clairezz/golab/hpool/hpool"
	"github.com/clairezz/golab/util"
)

const (
	JP_CNT = 10
)

func init() {
	util.InitLogger()
}

// 用户是怎么使用这个带负载均衡的pool的？
// 方案1 ：用户直接把任务交给 pool
// 方案2 ：用户从 pool 获得一个 worker （heap）, 用户把任务交给 worker （线程安全的map）
// 以下采用方案1

func main() {
	var maxNum = 100

	rand.Seed(time.Now().UnixNano())
	p := hpool.NewHpool()

	wg := sync.WaitGroup{}
	// job producers
	for i := 0; i < JP_CNT; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for rand.Intn(maxNum) != 0 {
				ch := make(chan string, 1) // buffer 为 1
				job := &hjob.Hjob{
					Handler: SayHi,
					Data:    []byte(strconv.Itoa(idx)),
					RespCh:  ch,
				}
				key, err := p.PushJob(job)
				if err != nil {
					seelog.Debugf(err.Error())
					continue
				}
				seelog.Debugf("client %d push job, got key %s\n", idx, key)

				resp := <-ch // 阻塞等结果
				seelog.Debugf("client %d got resp: %s\n", idx, resp)
			}
		}(i)
	}

	wg.Wait()

}

func SayHi(data []byte, respCh chan string) {
	respCh <- fmt.Sprintf("hello %s", string(data))
}
