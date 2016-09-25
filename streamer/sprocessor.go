package streamer

import (
	"sync"
)

/**
The processor interface shall be implemented by custom logic that processes messages in some way (e.g. transformation,
filtering, publishing, etc).
 */
type Processor interface {
	Execute(name string, cfg Config, input Message, out* chan Message)
}

/**
The base processor is the orchestrator of the processor execution and handles concurrency aspects.
 */
type BaseProcessor struct {
	Delegate Processor
	Balancer Demux
}

/**
The base execute method starts multiple routines for a processor depending on the balancer configuration
(e.g. parallelism hint) and waits for them to be complete.
 */
func (processor *BaseProcessor) Execute(name string, cfg Config, input <- chan Message) <- chan Message {
	var wg sync.WaitGroup
	numTasks := processor.Balancer.GetFanOut()
	wg.Add(numTasks)

	out := make(chan Message)

	work := func(taskId int, inputStream <- chan Message) {
		for message := range inputStream {
			processor.Delegate.Execute(name, cfg, message, &out)
		}
		wg.Done()
	}

	go func() {
		processor.Balancer.Run(input)
		for i := 0; i < numTasks; i++ {
			go work(i, processor.Balancer.GetOut(i))
		}
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
