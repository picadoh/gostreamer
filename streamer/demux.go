package streamer

/**
The demux interface provides the signature for message demultiplexers (e.g. run, get output channel for a given index
and get the number of output channels).
 */
type Demux interface {
	Run(input <- chan Message)
	GetOut(index int) <- chan Message
	GetFanOut() int
}
