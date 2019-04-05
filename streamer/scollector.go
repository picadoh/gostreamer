package streamer

import (
	"sync"
)

/**
The base collector is the orchestrator of the collector execution and handles the concurrency aspects.
*/
type Collector struct {
	name    string
	cfg     Config
	collect CollectFunction
}

type CollectFunction func(name string, cfg Config, out chan Message)

/**
The base execute method starts the delegate collector inside a routine and waits for it to finish.
*/
func (collector *Collector) Execute() <-chan Message {
	var wg sync.WaitGroup
	wg.Add(1)

	out := make(chan Message)

	go func() {
		collector.collect(collector.name, collector.cfg, out)
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

/**
Creates a new collector.
*/
func NewCollector(name string, cfg Config, collect CollectFunction) Collector {
	return Collector{name: name, cfg: cfg, collect: collect}
}
