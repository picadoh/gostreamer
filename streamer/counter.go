package streamer

import (
	"sync"
	"fmt"
)

/**
The counter represents a thread-safe key/value counting structure.
 */
type Counter struct {
	Mutex sync.RWMutex
	value map[string]int
}

func (counter *Counter) Increment(key string) int {
	counter.Mutex.Lock()
	defer counter.Mutex.Unlock()
	counter.value[key]++
	return counter.value[key]
}

func (counter *Counter) GetValue(key string) int {
	counter.Mutex.Lock()
	defer counter.Mutex.Unlock()
	return counter.value[key]
}

func (counter *Counter) ToString() string {
	return fmt.Sprintf("%s", counter.value)
}

func NewCounter() *Counter {
	return &Counter{value:make(map[string]int)}
}