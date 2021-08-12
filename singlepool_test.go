package async_utils

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestPool(t *testing.T) {
	pool := NewSinglePool(30, func() {
		fmt.Println("over")
	})

	var mu sync.Mutex
	var ic int

	for i := 0; i < 100; i++ {
		i2 := i
		pool.Send(func() error {
			mu.Lock()
			defer mu.Unlock()

			ic += 10
			if i2 == 60 || i2 == 35 {
				return errors.New("i2 == 60 || i2 == 35")
			}
			//i2 = i2
			return nil
		})
	}

	pool.Over()
	err := pool.Error()
	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Println(ic)
}
