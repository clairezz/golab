package main

import (
	"fmt"
	"time"

	"github.com/cihub/seelog"

	"golab/delayedTask/ring"
	"golab/util"
)


func init() {
	util.InitLogger()
}

func main() {
	defer seelog.Flush() // 进程结束前将内存中的日志写入文件
	seelog.Debugf("hello world")

	// 创建一个时间间隔为DURATION秒的环（环的一圈为1小时，环的当前指针每DURATION移动一次）
	l := ringr.NewLoop()
	seelog.Debugf("l.Len: %d", l.Len())

	// 添加延时任务
	for i :=  1; i < 23; i++ {
		l.AddTask(*ringr.NewTask(i,fmt.Sprintf("hello %2d", i))) // 任务是：i秒之后发消息
	}
	l.Status() // for debug：查看环的状态

	// 开始运行
	l.Run()

	time.Sleep(2*time.Second)

	// 添加延时任务
	for i :=  1; i < 23; i = i + 2 {
		l.AddTask(*ringr.NewTask(i,fmt.Sprintf("world %d", i))) // 任务是：i秒之后发消息
	}
	l.Status() // for debug：查看环的状态

	time.Sleep(60*time.Second)
	l.Stop()
}
