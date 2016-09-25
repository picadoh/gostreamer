package streamer

import (
	"sync"
)

/**
The counter represents a thread-safe key/value counting structure.
 */
type Counter struct {
	Mutex sync.RWMutex
	Count map[string]int
}

func (counter *Counter) Increment(key string) int {
	counter.Mutex.Lock()
	defer counter.Mutex.Unlock()
	counter.Count[key]++
	return counter.Count[key]
}

func (counter *Counter) GetValue(key string) int {
	counter.Mutex.Lock()
	defer counter.Mutex.Unlock()
	return counter.Count[key]
}

func NewCounter() *Counter {
	return &Counter{Count:make(map[string]int)}
}
