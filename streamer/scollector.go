package streamer

import (
	"sync"
)

/*
The collector interface shall be implemented with custom logic that collects information from a source.
 */
type Collector interface {
	Execute(name string, cfg Config, out* chan Message)
}

/**
The base collector is the orchestrator of the collector execution and handles the concurrency aspects.
 */
type BaseCollector struct {
	Delegate Collector
	Next Processor
}

/**
The base execute method starts the delegate collector inside a routine and waits for it to finish.
 */
func (collector *BaseCollector) Execute(name string, cfg Config) <- chan Message {
	var wg sync.WaitGroup
	wg.Add(1)

	out := make(chan Message)

	work := func(name string) {
		collector.Delegate.Execute(name, cfg, &out)
		wg.Done()
	}

	go func() {
		go work(name)
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
