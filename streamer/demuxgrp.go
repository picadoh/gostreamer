package streamer

import (
	"sync"
)

/**
The group demux splits an input channel into multiple output channels based on a given message key,
by retrieving the value and hashing it on a module of the number of output channels.
 */
type GroupDemux struct {
	Demux
	out      [] chan Message // Output channels
	groupKey string          // group distribution key
}

/**
Builds a new group demux based on the number of specified output channels and the key to be used in the group.
 */
func NewGroupDemux(nOutChannels int, groupKey string) *GroupDemux {
	demux := &GroupDemux{groupKey:groupKey}
	demux.out = make([]chan Message, nOutChannels)
	for i := 0; i < nOutChannels; i++ {
		demux.out[i] = make(chan Message)
	}
	return demux
}

/**
Executes the demultiplex function by assigning an index to a message based on its value (all messages with the same
value will be assigned the same index).
 */
func (demux *GroupDemux) Run(input <- chan Message) {
	nchannels := len(demux.out);

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for message := range input {
			value := message.Get(demux.groupKey)
			boundhash := int(Hash(value, nchannels))

			// Emit message
			demux.out[boundhash] <- message
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
func (demux *GroupDemux) GetOut(index int) <- chan Message {
	return demux.out[index]
}

/**
Gets the number of output channels.
 */
func (demux *GroupDemux) GetFanOut() int {
	return len(demux.out);
}
