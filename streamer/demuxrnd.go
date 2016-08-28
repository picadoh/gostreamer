package streamer

import (
	"sync"
	"math/rand"
)

type RandomDemux struct {
	Demux
	out [] chan Message // Output channels
}

func NewRandomDemux(nOutChannels int) *RandomDemux {
	demux := &RandomDemux{}
	demux.out = make([]chan Message, nOutChannels)
	for i := 0; i < nOutChannels; i++ {
		demux.out[i] = make(chan Message)
	}
	return demux
}

func (demux *RandomDemux) Run(input <- chan Message) {
	nchannels := len(demux.out);

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for message := range input {
			index := rand.Intn(nchannels)

			// Emit message
			demux.out[index] <- message
			//log.Printf("[randomdemux] Assigned hash %d for %s\n", index, message)
		}

		wg.Done()
	}()

	go func() {
		wg.Wait()

		for i := 0; i < nchannels; i++ {
			//log.Printf("[groupdemux] out[%d]:%d", i, len(demux.out[i]))
			close(demux.out[i])
		}
	}()
}

func (demux *RandomDemux) GetOut(index int) <- chan Message {
	return demux.out[index]
}

func (demux *RandomDemux) GetFanOut() int {
	return len(demux.out);
}
