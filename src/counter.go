package main

import (
	"sync"
)

type Counter struct {
	mutex sync.RWMutex
	count map[string]int
}

func (counter *Counter) Increment(key string) int {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()
	counter.count[key]++
	return counter.count[key]
}

func (counter *Counter) GetValue(key string) int {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()
	return counter.count[key]
}

func NewCounter() *Counter {
	return &Counter{count:make(map[string]int)}
}