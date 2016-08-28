package streamer

import (
	"sync"
)

type GroupDemux struct {
	Demux
	out      [] chan Message // Output channels
	groupKey string          // group distribution key
}

func NewGroupDemux(nOutChannels int, groupKey string) *GroupDemux {
	demux := &GroupDemux{groupKey:groupKey}
	demux.out = make([]chan Message, nOutChannels)
	for i := 0; i < nOutChannels; i++ {
		demux.out[i] = make(chan Message)
	}
	return demux
}

func (demux *GroupDemux) Run(input <- chan Message) {
	nchannels := len(demux.out);

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for message := range input {
			value := message.Get(demux.groupKey)
			boundhash := int(hash(value, nchannels))

			// Emit message
			demux.out[boundhash] <- message
			//log.Printf("[groupdemux] Assigned hash %d for %s/%s\n", boundhash, demux.groupKey, message)
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

func (demux *GroupDemux) GetOut(index int) <- chan Message {
	return demux.out[index]
}

func (demux *GroupDemux) GetFanOut() int {
	return len(demux.out);
}
