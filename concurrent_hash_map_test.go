package async_utils

import (
	"fmt"
	"log"
	"testing"
)

func TestConcurrentHashMap(t *testing.T) {
	hashMap := NewConcurrentHashMap()
	for i := 0; i < 100; i++ {
		hashMap.Insert(fmt.Sprintf("key:%d", i), i)
	}

	for i := 0; i < 100; i++ {
		get, err := hashMap.Get(fmt.Sprintf("key:%d", i))
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(get)
	}

	for i := 0; i < 100; i++ {
		hashMap.Insert(fmt.Sprintf("key:%d", i), i)
	}

	hashMap.Iteration(func(key string, val interface{}) {
		fmt.Printf("key: (%s) val: %v \n", key, val)
	})
}
