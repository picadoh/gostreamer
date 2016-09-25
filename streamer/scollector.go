package streamer

import (
	"sync"
)

/**
The base collector is the orchestrator of the collector execution and handles the concurrency aspects.
 */
type Collector struct {
	Name    string
	Cfg     Config
	Collect CollectFunction
}

type CollectFunction func(name string, cfg Config, out* chan Message)

/**
The base execute method starts the delegate collector inside a routine and waits for it to finish.
 */
func (collector *Collector) Execute() <- chan Message {
	var wg sync.WaitGroup
	wg.Add(1)

	out := make(chan Message)

	work := func(name string) {
		collector.Collect(name, collector.Cfg, &out)
		wg.Done()
	}

	go func() {
		go work(collector.Name)
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func NewCollector(name string, cfg Config, collect CollectFunction) Collector {
	return Collector{Name:name,Cfg:cfg,Collect:collect}
}
