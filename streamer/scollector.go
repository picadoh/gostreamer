package streamer

import (
	"sync"
)

type CollectorFunction func(cfg Config, out* chan Message)

func SCollector(cfg Config, name string, collector CollectorFunction) <- chan Message {
	var wg sync.WaitGroup
	wg.Add(1)

	out := make(chan Message)

	work := func(name string) {
		//log.Printf("[%s] starting\n", name)
		collector(cfg, &out)
		//log.Printf("[%s] ending\n", name)

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
