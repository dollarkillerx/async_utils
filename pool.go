package async_utils

import (
	"fmt"
	"sync"
)

type PoolFunc func()

type easyPool struct {
	limit chan bool
	pool  chan PoolFunc
	close PoolFunc
}

// 方法池
func NewPoolFunc(size int, close PoolFunc) *easyPool {
	pool := &easyPool{
		limit: make(chan bool, size),
		pool:  make(chan PoolFunc, 10),
		close: close,
	}
	go pool.core()
	return pool
}

//  下发任务
func (e *easyPool) Send(fn PoolFunc) {
	e.pool <- fn
}

// 下发完毕信号  (任务下发完毕时调用)
func (e *easyPool) Over() {
	close(e.pool)
}

func (e *easyPool) core() {
	wg := sync.WaitGroup{}
loop:
	for {
		select {
		case fn, over := <-e.pool:
			if !over {
				break loop
			}
			e.limit <- true
			wg.Add(1)
			go func(fn PoolFunc) {
				defer func() {
					if err := recover(); err != nil {
						PrintStack()
						fmt.Println("Recover Err: ", err)
					}
				}() // 我怕你们乱写逻辑 把系统弄炸了

				defer func() {
					<-e.limit
					wg.Done()
				}()
				fn()
			}(fn)
		}
	}

	wg.Wait()
	e.close()
}

// 再来一个简单任务池
type SimpleTask struct {
	funcs chan func()
	ove   chan bool
	wg    sync.WaitGroup
}

func NewSimpleTask(limit int) *SimpleTask {
	sim := SimpleTask{
		funcs: make(chan func(), limit),
		ove:   make(chan bool),
	}
	go sim.core()
	return &sim
}

//func (s *SimpleTask) Over() {
//	close(s.funcs)
//}

func (s *SimpleTask) AddTask(fn func()) {
	s.funcs <- fn
}

func (s *SimpleTask) Wait() {
	close(s.funcs)
	<-s.ove
}

func (s *SimpleTask) core() {
loop:
	for {
		select {
		case fn, exit := <-s.funcs:
			if !exit {
				break loop
			}

			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				fn()
			}()
		}
	}

	s.wg.Wait()
	s.ove <- true
}
