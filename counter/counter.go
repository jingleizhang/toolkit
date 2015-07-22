package counter

import (
	"sync"
	"time"
)

/*
 * A simple thread safe counter.
 */
type Counter struct {
	data      map[string]int64
	StartTime string
	lock      *sync.RWMutex
}

func NewCounter() *Counter {
	return &Counter{
		data:      make(map[string]int64),
		StartTime: time.Now().Format("2006/01/02 15:04:05"),
		lock:      &sync.RWMutex{},
	}
}

func (self *Counter) Incr(key string, val int64) {
	self.lock.Lock()
	defer self.lock.Unlock()
	_, ok := self.data[key]
	if !ok {
		self.data[key] = val
	} else {
		self.data[key] += val
	}
}

func (self *Counter) Get(key string) int64 {
	self.lock.RLock()
	defer self.lock.RUnlock()
	val, ok := self.data[key]
	if !ok {
		return 0
	}
	return val
}

func (self *Counter) Set(key string, val int64) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.data[key] = val
}

func (self *Counter) Stat() map[string]int64 {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.data
}

func (self *Counter) Data() map[string]int64 {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.data
}

func (self *Counter) Reset() {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.data = make(map[string]int64)
	self.StartTime = time.Now().Format("2006/01/02 15:04:05")
}
