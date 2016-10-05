package streamer

import "sync"

/**
The demux interface provides the signature for message demultiplexers (e.g. run, get output channel for a given index
and get the number of output channels).
 */
type ChannelDemux interface {
	Execute(input <- chan Message)
	Output(index int) <- chan Message
	FanOut() int
}

/**
Implementation that makes use of an index function to handle the calculation of the index
for a given message.
 */
type IndexedChannelDemux struct {
	ChannelDemux
	out   [] chan Message
	index IndexFunction
}

type IndexFunction func(nchannels int, element Message) int

/**
Executes the demultiplex function from context which assigns an index to a message.
 */
func (demux *IndexedChannelDemux) Execute(input <- chan Message) {
	nchannels := len(demux.out);

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for message := range input {
			// assign index
			index := demux.index(nchannels, message)

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
func (demux *IndexedChannelDemux) Output(index int) <- chan Message {
	return demux.out[index]
}

/**
Gets the number of output channels.
 */
func (demux *IndexedChannelDemux) FanOut() int {
	return len(demux.out);
}

/**
Builds a new group demux based on the number of specified output channels and the key to be used in the group.
 */
func NewIndexedChannelDemux(fanOut int, index IndexFunction) ChannelDemux {
	demux := &IndexedChannelDemux{index: index}
	demux.out = make([]chan Message, fanOut)
	for i := 0; i < fanOut; i++ {
		demux.out[i] = make(chan Message)
	}
	return demux
}