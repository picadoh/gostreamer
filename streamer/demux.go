package streamer

type Demux interface {
	Run(input <- chan Message)
	GetOut(index int) <- chan Message
	GetFanOut() int
}
