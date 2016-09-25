package streamer

import "math/rand"

/**
Demux context interface provides the signature interface for executing the demultiplexing
under a context.
 */
type DemuxContext interface {
	Execute(nOutChannels int, input Message) int
}

/**
Group demux context implements the demultiplexing function through the use of a group
key that is used to retrieve the value and perform a hash module over it to retrieve the index.
 */
type GroupDemuxContext struct {
	DemuxContext
	group string
}

func (fun *GroupDemuxContext) Execute(nOutChannels int, input Message) int {
	value := input.Get(fun.group)
	return int(Hash(value, nOutChannels))
}

func NewGroupDemuxCtx(group string) DemuxContext {
	return &GroupDemuxContext{group:group}
}

/*
Random demux context implements the demultiplexing function that generates the index
randomly.
 */
type RandomDemuxContext struct {
	DemuxContext
	group string
}

func (fun *RandomDemuxContext) Execute(nChannels int, input Message) int {
	return rand.Intn(nChannels)
}

func NewRandomDemuxCtx() DemuxContext {
	return &RandomDemuxContext{}
}
