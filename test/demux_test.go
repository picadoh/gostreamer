package streamer_test

import (
	"github.com/picadoh/gostreamer/streamer"
	"testing"
)

func TestDemuxContextFanOut(t *testing.T) {
	victim := streamer.NewDemux(5, NewMockDemuxCtx())

	if (victim.GetFanOut() != 5) {
		t.Errorf("Expected fan out 5, found %d", victim.GetFanOut())
	}
}

func TestDemuxExecute(t *testing.T) {
	// init the victim
	victim := streamer.NewDemux(3, NewMockDemuxCtx())

	// prepare the scenario
	input := make(chan streamer.Message)

	msg := streamer.NewMessage()
	msg.Put("testkey", "testvalue")

	go func() {
		input <- msg
	}()

	victim.Execute(input)
	output := <-victim.GetOut(1)

	if (output == nil) {
		t.Errorf("Expected nothing, found %s", output)
	}

	if (output.Get("testkey") != "testvalue") {
		t.Errorf("Expected testvalue, found %s", output.Get("testkey"))
	}
}

type MockDemuxContext struct {
	streamer.DemuxContext
}

func (fun *MockDemuxContext) Execute(nChannels int, input streamer.Message) int {
	return 1
}

func NewMockDemuxCtx() streamer.DemuxContext {
	return &MockDemuxContext{}
}
