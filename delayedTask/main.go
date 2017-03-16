package main

import (
	"fmt"
	"time"
	"golab/delayedTask/ring"
)

func main() {
	fmt.Println("hello world")
	l := ring.NewLoop(10)
	l.Run()
	for i :=  1; i < 100; i++ {
		l.AddTsk(*ring.NewTask(i,fmt.Sprintf("hello %d\n", i)))
	}
	time.Sleep(60*time.Second)
}
