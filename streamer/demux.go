package streamer

import (
	"sync"
)

/**
The demux interface provides the signature for message demultiplexers (e.g. run, get output channel for a given index
and get the number of output channels).
 */
type Demux interface {
	Execute(input <- chan Message)
	GetOut(index int) <- chan Message
	GetFanOut() int
}

/**
Implementation that makes use of a demux context to handle the calculation of the index
for a given message.
 */
type DemuxImpl struct {
	Demux
	out [] chan Message // Output channels
	ctx DemuxContext
}

/**
Executes the demultiplex function from context which assigns an index to a message.
 */
func (demux *DemuxImpl) Execute(input <- chan Message) {

	nchannels := len(demux.out);

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for message := range input {
			// assign index
			index := demux.ctx.Execute(nchannels, message)

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
func (demux *DemuxImpl) GetOut(index int) <- chan Message {
	return demux.out[index]
}

/**
Gets the number of output channels.
 */
func (demux *DemuxImpl) GetFanOut() int {
	return len(demux.out);
}

/**
Builds a new group demux based on the number of specified output channels and the key to be used in the group.
 */
func NewDemux(nOutChannels int, context DemuxContext) *DemuxImpl {
	demux := &DemuxImpl{ctx:context}
	demux.out = make([]chan Message, nOutChannels)
	for i := 0; i < nOutChannels; i++ {
		demux.out[i] = make(chan Message)
	}
	return demux
}