package streamer

import (
	"sync"
	"math/rand"
)

/**
The random demux splits an input channel into multiple output channels
based on a random index between 0 and the number of channels.
 */
type RandomDemux struct {
	Demux
	out [] chan Message // Output channels
}

/**
Builds a new random demux with a specified number of output channels.
 */
func NewRandomDemux(nOutChannels int) *RandomDemux {
	demux := &RandomDemux{}
	demux.out = make([]chan Message, nOutChannels)
	for i := 0; i < nOutChannels; i++ {
		demux.out[i] = make(chan Message)
	}
	return demux
}

/**
Executes the demultiplex function that generates the output randomly.
 */
func (demux *RandomDemux) Run(input <- chan Message) {
	nchannels := len(demux.out);

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for message := range input {
			index := rand.Intn(nchannels)
			// Emit message
			demux.out[index] <- message
		}

		wg.Done()
	}()

	go func() {
		wg.Wait()
		for i := 0; i < nchannels; i++ {
			close(demux.out[i])
		}
	}()
}

/**
Gets the output channel for a given index.
 */
func (demux *RandomDemux) GetOut(index int) <- chan Message {
	return demux.out[index]
}

/**
Gets the number of output channels.
 */
func (demux *RandomDemux) GetFanOut() int {
	return len(demux.out);
}
