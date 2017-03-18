package ccp

import (
	"time"
	"sync"
	"math/rand"
	"fmt"
)

// Senders2Receiver 模拟了多个 sender， 一个 receiver 的应用场景下 channel 关闭方式：增加一个signal channel，buffer 为 1 ，唯一的receiver关闭
// dataChan 不用关闭了么？
func Senders2Receiver() {
	rand.Seed(time.Now().UnixNano())
	var dataChanBufCnt = 5 // channel的缓冲区大小
	var maxNum = 10
	var senderCnt = 10

	dataChan := make(chan int, dataChanBufCnt)
	doneChan := make(chan struct{}) // signal channel
	wg := sync.WaitGroup{}

	// senders
	index := 0
	for i := 0; i < senderCnt; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for {
				data := rand.Intn(maxNum)
				select {
				case <- doneChan: // sender 收到关闭信号 退出
					fmt.Printf("receiver %d exit.\n", idx)
					return
				case dataChan <- data:
				}
			}
		}(index)
		index++
	}

	// receiver
	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range dataChan {
			fmt.Printf("receive data %d\n", data)
			if data == maxNum - 1 {
				fmt.Printf("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n")
				close(doneChan) // 通知sender退出
				return
			}
		}
	}()

	wg.Wait() // 等待 所以有goroutine退出
}
