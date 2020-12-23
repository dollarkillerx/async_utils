package async_utils

import (
	"fmt"
	"sync"
)

type LockMap struct {
	mu map[string]*sync.Mutex
}

func NewLockMap() *LockMap {
	return &LockMap{
		mu: map[string]*sync.Mutex{},
	}
}

type RWLockMap struct {
	mu map[string]*sync.RWMutex
}

func NewRWLockMap() *RWLockMap {
	return &RWLockMap{
		mu: map[string]*sync.RWMutex{},
	}
}

func (l *LockMap) Lock(key string) {
	mutex, ex := l.mu[key]
	if !ex {
		l.mu[key] = &sync.Mutex{}
		mutex = l.mu[key]
	}

	mutex.Lock()
}

func (l *LockMap) UnLock(key string) error {
	mutex, ex := l.mu[key]
	if !ex {
		return fmt.Errorf("Unlock of unlocked Mutex")
	}

	mutex.Unlock()
	return nil
}

func (l *RWLockMap) Lock(key string) {
	mutex, ex := l.mu[key]
	if !ex {
		l.mu[key] = &sync.RWMutex{}
		mutex = l.mu[key]
	}

	mutex.Lock()
}

func (l *RWLockMap) UnLock(key string) error {
	mutex, ex := l.mu[key]
	if !ex {
		return fmt.Errorf("Unlock of unlocked RWMutex")
	}

	mutex.Unlock()
	return nil
}

func (l *RWLockMap) RLock(key string) {
	mutex, ex := l.mu[key]
	if !ex {
		l.mu[key] = &sync.RWMutex{}
		mutex = l.mu[key]
	}

	mutex.RLock()
}

func (l *RWLockMap) RUnlock(key string) error {
	mutex, ex := l.mu[key]
	if !ex {
		return fmt.Errorf("Unlock of unlocked RWMutex")
	}

	mutex.RUnlock()
	return nil
}
