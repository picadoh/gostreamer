package streamer_test

import (
	"github.com/picadoh/gostreamer/streamer"
	"testing"
	"log"
)

func TestProcessData(t *testing.T) {
	input := make(chan streamer.Message, 1)
	defer close(input)

	demuxout := make(chan streamer.Message)
	defer close(demuxout)

	msg := streamer.NewMessage()
	msg.Put("testkey", "testvalue")
	input <- msg

	cfg := streamer.NewPropertiesConfig()

	victim := streamer.NewProcessor("x", cfg, MockProcess, NewMockDemux(demuxout))

	victim.Execute(input)

	dmuxMsg := <-demuxout
	if dmuxMsg == nil {
		t.Error("Expected output, found nothing")
	}

	if dmuxMsg.Get("x") != "y" {
		t.Error("Expected y, found ", dmuxMsg.Get("x"))
	}
}

func MockProcess(name string, cfg streamer.Config, input streamer.Message, out chan streamer.Message) {
	log.Printf("Executing mocked process %s with config %s", name, cfg.ToString())
	out <- input
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

func (demux MockDemux) Execute(input <- chan streamer.Message) {
	msg := streamer.NewMessage()
	msg.Put("x", "y")
	demux.out[0] <- msg
}

func (demux MockDemux) GetOut(index int) <- chan streamer.Message {
	return demux.out[index]
}

func (demux MockDemux) GetFanOut() int {
	return len(demux.out);
}
