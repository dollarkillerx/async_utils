package async_utils

import (
	"fmt"
	"sync"

	"go.uber.org/atomic"
)

type PoolFunc func()

type EasyPool struct {
	limit chan struct{}
	pool  chan PoolFunc
	close PoolFunc
}

// NewPoolFunc 方法池
func NewPoolFunc(size int, close PoolFunc) *EasyPool {
	pool := &EasyPool{
		limit: make(chan struct{}, size),
		pool:  make(chan PoolFunc, size),
		close: close,
	}
	go pool.core()
	return pool
}

// Send 下发任务
func (e *EasyPool) Send(fn PoolFunc) {
	e.pool <- fn
}

// Over 下发完毕信号  (任务下发完毕时调用)
func (e *EasyPool) Over() {
	close(e.pool)
}

func (e *EasyPool) core() {
	wg := sync.WaitGroup{}
loop:
	for {
		select {
		case fn, over := <-e.pool:
			if !over {
				break loop
			}
			e.limit <- struct{}{}
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

type SinglePoolFunc func() error

// SinglePool 可停止的pool
type SinglePool struct {
	pool      chan SinglePoolFunc
	close     PoolFunc
	limit     chan struct{}
	err       error
	rmu       sync.Mutex
	stop      atomic.Bool // stop  控制错误关闭
	closeChan chan struct{}
}

func NewSinglePool(size int, close PoolFunc) *SinglePool {
	pool := &SinglePool{
		limit:     make(chan struct{}, size),
		pool:      make(chan SinglePoolFunc, size),
		close:     close,
		closeChan: make(chan struct{}),
	}
	go pool.core()
	return pool
}

// Send 下发任务
func (s *SinglePool) Send(fn SinglePoolFunc) {
	if !s.stop.Load() {
		s.pool <- fn
	}
}

// Over 下发完毕信号  (任务下发完毕时调用)
func (s *SinglePool) Over() {
	close(s.pool)
}

func (s *SinglePool) core() {
	wg := sync.WaitGroup{}
loop:
	for {
		select {
		case fn, over := <-s.pool:
			if !over {
				break loop
			}

			if s.stop.Load() {
				break loop
			}

			s.limit <- struct{}{}
			wg.Add(1)
			go func(fn SinglePoolFunc) {
				defer func() {
					if err := recover(); err != nil {
						PrintStack()
						fmt.Println("Recover Err: ", err)
					} // 我怕你们乱写逻辑 把系统弄炸了

					<-s.limit
					wg.Done()
				}()

				err := fn()
				if err != nil {
					s.rmu.Lock()
					s.err = err
					s.stop.Store(true)
					s.rmu.Unlock()
				}
			}(fn)
		}
	}

	// 判断是否是因为错误而停止
	if !s.stop.Load() {
		wg.Wait()
	}

	close(s.closeChan)
	s.close()
}

func (s *SinglePool) Error() error {
	<-s.closeChan
	return s.err
}

// Greedy Pool
type greedyPool struct {
	pool chan PoolFunc
	over chan struct{}

	close PoolFunc
}

// NewGreedyPool: return greed pool
func NewGreedyPool(size int, close PoolFunc) *greedyPool {
	pool := &greedyPool{
		pool:  make(chan PoolFunc, size),
		over:  make(chan struct{}),
		close: close,
	}

	go pool.scheduler(size)
	return pool
}

// Over: End of task issuance
func (g *greedyPool) Over() {
	close(g.over)
}

// Send: Issue a task
func (g *greedyPool) Send(fn PoolFunc) {
	g.pool <- fn
}

func (g *greedyPool) scheduler(size int) {
	var wg sync.WaitGroup
	for i := 0; i < size; i++ {
		wg.Add(1)
		go g.task(&wg)
	}

	wg.Wait()
	g.close()
}

func (g *greedyPool) task(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()

		// recover
		if err := recover(); err != nil {
			PrintStack()
			fmt.Println("Recover Err: ", err)
		}
	}()

loop:
	for {
		select {
		case i, ex := <-g.pool:
			if !ex {
				break loop
			}
			i()
		case <-g.over:
			break loop
		}
	}
}

// 再来一个简单任务池
type SimpleTask struct {
	funcs chan func()
	ove   chan struct{}
	wg    sync.WaitGroup
}

func NewSimpleTask(limit int) *SimpleTask {
	sim := SimpleTask{
		funcs: make(chan func(), limit),
		ove:   make(chan struct{}),
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
	s.ove <- struct{}{}
}
