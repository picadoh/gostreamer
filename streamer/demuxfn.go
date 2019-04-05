package streamer

import "math/rand"

/**
Group index function implements the demultiplexing function through the use of a group
key that is used to retrieve the value and perform a hash module over it to retrieve the index.
*/
type GroupDemux struct {
	group string
}

func (fun *GroupDemux) GroupIndex(nOutChannels int, input Message) int {
	value := input.(Message).Get(fun.group)
	return int(Hash(value, nOutChannels))
}

func NewGroupDemux(key string) *GroupDemux {
	return &GroupDemux{group: key}
}

/**
Random index function implements demultiplexing by randomly assigning an index, independently
 of the input message.
*/
func RandomIndex(fanOut int, input Message) int {
	return rand.Intn(fanOut)
}
