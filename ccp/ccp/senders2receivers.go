package ccp

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// TODO 死锁
// Senders2Receivers 模拟了多个 sender, 多个 receiver 的场景下关闭 channel 的方式： 由 mediator（调停者）来通知 senders 和 receivers 退出
func Senders2Receivers() {
	rand.Seed(time.Now().UnixNano())

	var senderCnt = 20
	var receiverCnt = 12
	var maxNum = 10000
	var dataBuf = 30

	dataCh := make(chan int, dataBuf)  //
	doneCh := make(chan struct{})      // notify 所有 sender 和 receiver 退出
	toStopCh := make(chan struct{}, 1) // 通知 mediator 发退出消息

	wg := sync.WaitGroup{}

	// mediator
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-toStopCh:
				fmt.Printf("@@@@@@@@@@@@@@@@\n")
				close(doneCh)
				return
			default:
				// do nothing
			}
		}
	}()

	// sender
	for i := 0; i < senderCnt; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for {
				select {
				case <-doneCh:
					fmt.Printf("sender %d exit\n", idx)
					return
				default:
					data := rand.Intn(maxNum)
					if data == maxNum-1 {
						select {
						case toStopCh <- struct{}{}:
							fmt.Printf("sender %d stops the process.\n", idx)
						default:
							// do nothing
						}
					} else {
						dataCh <- data
					}
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
				// 尽快发现关闭信号
				select {
				case <-doneCh:
					fmt.Printf("sender %d exit\n", idx)
					return
				default:
					// do noting
				}

				select {
				case <-doneCh:
					fmt.Printf("sender %d exit\n", idx)
					return
				case data := <-dataCh:
					fmt.Printf("receiver %d receive data %d\n", idx, data)
					if data == maxNum-1 {
						select {
						case toStopCh <- struct{}{}:
							fmt.Printf("receiver %d stops the process\n", idx)
						default:
						}
					}
				}
			}
		}(i)
	}

	wg.Wait()
}
