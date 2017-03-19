package main

import (
	"golab/hpool/hpool"
	"fmt"
	"golab/hpool/hjob"
	"golab/hpool/hworker"
	"time"
	"math/rand"
)

const (
	JP_CNT = 10
)

// 用户是怎么使用这个带负载均衡的pool的？
// 方案1 ：用户直接把任务交给 pool
// 方案2 ：用户从 pool 获得一个 worker （heap）, 用户把任务交给 worker （线程安全的map）
// 以下采用方案2

func main() {
	var maxNum = 100

	rand.Seed(time.Now().UnixNano())
	p := hpool.NewHpool()

	// job producers
	for i := 0; i < JP_CNT; i++ {
		go func() {
			for rand.Intn(maxNum) != 0 {
				j := &hjob.Hjob{
					Handler:SayHi,
				}
				item := p.Pop()
				if item == nil {
					fmt.Println("86868686868")
					continue
				}
				worker := item.(*hworker.Hworker)
				worker.AddJob(j)
				p.Push(worker) // ！！！用完记得还回去
			}
		}()
	}

	time.Sleep(60e9)

}

func SayHi() {
	fmt.Println("Hello World")
}