package streamer

import (
	"sync"
)

type ProcessorFunction func(cfg Config, input Message, out*chan Message)

func SProcessor(cfg Config, name string, demux Demux, input <- chan Message, processor ProcessorFunction) <- chan Message {
	var wg sync.WaitGroup
	numTasks := demux.GetFanOut()
	wg.Add(numTasks)

	out := make(chan Message)

	work := func(taskId int, inputStream <- chan Message) {
		//log.Printf("[%s] starting task %d\n", name, taskId)

		for message := range inputStream {
			//log.Printf("[%s] Task %d picked up message %s\n", name, taskId, message)
			processor(cfg, message, &out)
		}

		//log.Printf("[%s] ending task %d\n", name, taskId)

		wg.Done()
	}

	go func() {
		demux.Run(input)
		for i := 0; i < numTasks; i++ {
			go work(i, demux.GetOut(i))
		}
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
