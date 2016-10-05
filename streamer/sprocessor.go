package streamer

import (
	"sync"
)

/**
The base processor is the orchestrator of the processor execution and handles concurrency aspects.
 */
type Processor struct {
	name    string
	cfg     Config
	process ProcessFunction
	demux   ChannelDemux
}

type ProcessFunction func(name string, cfg Config, input Message, out chan Message)

/**
The base execute method starts multiple routines for a processor depending on the balancer configuration
(e.g. parallelism hint) and waits for them to be complete.
 */
func (processor *Processor) Execute(input <- chan Message) <- chan Message {
	var wg sync.WaitGroup
	numTasks := processor.demux.FanOut()
	wg.Add(numTasks)

	out := make(chan Message)

	work := func(taskId int, inputStream <- chan Message) {
		for message := range inputStream {
			processor.process(processor.name, processor.cfg, message, out)
		}
		wg.Done()
	}

	go func() {
		processor.demux.Execute(input)
		for i := 0; i < numTasks; i++ {
			go work(i, processor.demux.Output(i))
		}
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

/**
Creates a new processor.
 */
func NewProcessor(name string, cfg Config, process ProcessFunction, demux ChannelDemux) Processor {
	return Processor{name:name, cfg:cfg, process:process, demux:demux}
}