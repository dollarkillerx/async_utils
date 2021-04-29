package async_utils

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestGreedyPool(t *testing.T) {
	ov := make(chan struct{})
	pool := NewGreedyPool(30, func() {
		close(ov)
	})

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println(runtime.NumGoroutine())
		}
	}()

	time.Sleep(time.Second)
	for i := 0; i < 99999999999; i++ {
		idx := i
		pool.Send(func() {
			idx += 36
		})
	}

	pool.Over()
	<-ov
	fmt.Println("Over")
}
