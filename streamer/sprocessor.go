package streamer

import (
	"sync"
)

/**
The base processor is the orchestrator of the processor execution and handles concurrency aspects.
 */
type Processor struct {
	Name string
	Cfg Config
	Process ProcessFunction
	Balancer Demux
}

type ProcessFunction func(name string, cfg Config, input Message, out* chan Message)

/**
The base execute method starts multiple routines for a processor depending on the balancer configuration
(e.g. parallelism hint) and waits for them to be complete.
 */
func (processor *Processor) Execute(input <- chan Message) <- chan Message {
	var wg sync.WaitGroup
	numTasks := processor.Balancer.GetFanOut()
	wg.Add(numTasks)

	out := make(chan Message)

	work := func(taskId int, inputStream <- chan Message) {
		for message := range inputStream {
			processor.Process(processor.Name, processor.Cfg, message, &out)
		}
		wg.Done()
	}

	go func() {
		processor.Balancer.Execute(input)
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

func NewProcessor(name string, cfg Config, process ProcessFunction, balancer Demux) Processor  {
	return Processor{Name:name,Cfg:cfg,Process:process,Balancer:balancer}
}