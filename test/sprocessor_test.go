package streamer_test

import (
	"github.com/picadoh/gostreamer/streamer"
	"testing"
)

func TestCollectData(t *testing.T) {

	input := make(chan streamer.Message)
	demuxout := make(chan streamer.Message)
	done := make(chan bool)

	go func() {
		msg := streamer.NewMessage()
		msg.Put("testkey", "testvalue")
		input <- msg
	}()

	<-input

	cfg := &streamer.PropertiesConfig{Properties:make(map[string]string)}

	victim := streamer.NewProcessor("x", cfg, MockProcess, NewMockDemux(demuxout))

	var output <- chan streamer.Message
	go func() {
		output = victim.Execute(input)
		done <- true
	}()

	<-done

	close(done)
	close(input)

	dmuxMsg := <-demuxout
	if dmuxMsg == nil {
		t.Error("Expected output, found nothing")
	}

	if dmuxMsg.Get("x") != "y" {
		t.Error("Expected y, found ", dmuxMsg.Get("x"))
	}

	close(demuxout)
}

func MockProcess(name string, cfg streamer.Config, input streamer.Message, out *chan streamer.Message) {
	*out <- input
}

type MockDemux struct {
	streamer.Demux
	out [] chan streamer.Message // Output channels
}

func NewMockDemux(output chan streamer.Message) streamer.Demux {
	demux := MockDemux{}
	demux.out = make([]chan streamer.Message, 1)
	demux.out[0] = output
	return demux
}

func (demux MockDemux) Run(input <- chan streamer.Message) {
	msg := streamer.NewMessage()
	msg.Put("x","y")
	demux.out[0] <- msg
}

func (demux MockDemux) GetOut(index int) <- chan streamer.Message {
	return demux.out[index]
}

func (demux MockDemux) GetFanOut() int {
	return len(demux.out);
}
