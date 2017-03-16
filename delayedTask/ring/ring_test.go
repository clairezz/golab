package ring

import (
	"testing"
	"fmt"
	"time"
)

func TestNewLoop(t *testing.T) {
	fmt.Println("hello world")
	l := NewLoop(10)
	l.Run()
	for i :=  1; i < 100; i++ {
		l.AddTsk(*NewTask(i,fmt.Sprintf("hello %d\n", i)))
	}
	time.Sleep(60*time.Second)
}
