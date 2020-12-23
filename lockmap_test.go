package async_utils

import (
	"fmt"
	"testing"
	"time"
)

func TestLookMap(t *testing.T) {
	lockMap := NewLockMap()
	lockMap.Lock("abc")

	go func() {
		fmt.Println("in")
		lockMap.Lock("abc")
		fmt.Println("get look")
	}()

	time.Sleep(3 * time.Second)
	lockMap.UnLock("abc")
	time.Sleep(30 * time.Second)
}
