package ccp

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Senders2Receivers 模拟了多个 sender, 多个 receiver 的场景下关闭 channel 的方式： 由 mediator（调停者）来通知 senders 和 receivers 退出
func Senders2Receivers() {
	rand.Seed(time.Now().UnixNano())

	var senderCnt = 1
	var receiverCnt = 1
	var maxNum = 1000
	var dataBuf = 1

	dataCh := make(chan int, dataBuf)  //
	doneCh := make(chan struct{})      // notify 所有 sender 和 receiver 退出
	toStopCh := make(chan struct{}, 1) // 通知 mediator 发退出消息, 注意 buffer 为1

	wg := sync.WaitGroup{}

	// mediator
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-toStopCh // 等待某个 sender 或 receiver 发送停止消息
		fmt.Println("@@@@@@@@@@@@@@@@")
		close(doneCh) // 通知所有 sender 和 receiver 退出
	}()

	// sender
	for i := 0; i < senderCnt; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for {
				data := rand.Intn(maxNum)
				if data == maxNum-1 { // ！！！ 先判断自己是否需要停止，如果不需要停止，再去检查 doneCh 信号和发送 daata
					select {
					case toStopCh <- struct{}{}:
					// 发送
						fmt.Printf("sender %d stops the process.\n", idx)
					default: // 已满，直接退出
					}
					return // 只要条件满足，一定会退出，如果是第一个绝对退出的，需要发送toStop消息通知 mediator 去通知所有其他 sender 或 receiver 退出
				}

				select {
				case <-doneCh: // 等待 mediator close doneCH 来通知所有协程退出
					fmt.Printf("sender %d exit\n", idx)
					return
				case dataCh <- data: // 否则，发送数据
					fmt.Printf("sender %d send data %d\n", idx, data)
				}
			}
		}(i)
	}

	// receiver
	for i := 0; i < receiverCnt; i++ {
		wg.Add(1)
		go func(idx int) {
			wg.Done()
			for {
				select {
				case <-doneCh:
					fmt.Printf("receiver %d exit\n", idx)
					return
				case data := <-dataCh:
					fmt.Printf("receiver %d received data %d\n", idx, data)
					if data == 0 {
						select{
						case toStopCh <- struct{}{}: // 发送
							fmt.Printf("sender %d stops the process.\n", idx)
						default: // 已满，直接退出
						}
						return // 只要条件满足，一定会退出，如果是第一个绝对退出的，需要发送toStop消息通知 mediator 去通知所有其他 sender 或 receiver 退出
					}
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("-------------DONE----------------")
	fmt.Printf("reminded data cnt: %d\n", len(dataCh))
}
