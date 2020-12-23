package async_utils

import "sync"

// 数据异步存储
type AsyncSlice struct {
	slice []interface{}
	mu    sync.Mutex
}

func NewAsyncSlice() *AsyncSlice {
	slice := AsyncSlice{
		slice: []interface{}{},
	}
	return &slice
}

func (a *AsyncSlice) Append(data interface{}) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.slice = append(a.slice, data)
}

func (a *AsyncSlice) Get() []interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.slice
}
