package ccp

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Sender2Receivers 模拟了一个 sender，多个 receiver 的场景下 channel 关闭方式
func Sender2Receivers() {
	rand.Seed(time.Now().UnixNano())
	var maxNum = 1000

	var dataChanBufCnt = 5 // channel的缓冲区大小

	dataChan := make(chan int, dataChanBufCnt)
	wg := sync.WaitGroup{}

	// receiver
	index := 0 // receiver的编号
	for rand.Intn(maxNum) != 0 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			// 从 channel 中接收数据直到channel
			for data := range dataChan {
				fmt.Printf("receiver %d receives data %d\n", idx, data)
			}
			fmt.Printf("*********receiver %d exits*********\n", idx)
		}(index)
		index++
	}

	// sender
	for {
		// 当随机数为 maxNum-1 时，sender 关闭 channel 并退出循环
		d := rand.Intn(maxNum)
		if d == maxNum-1 {
			close(dataChan)
			break
		} else {
			dataChan <- d
		}
	}
	// 等待接收者退出
	wg.Wait()
}
