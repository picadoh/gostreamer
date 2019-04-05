package streamer_test

import (
	"testing"

	"github.com/picadoh/gostreamer/streamer"
)

func TestDemuxContextFanOut(t *testing.T) {
	victim := streamer.NewIndexedChannelDemux(5, MockIndexFunction)

	if victim.FanOut() != 5 {
		t.Errorf("Expected fan out 5, found %d", victim.FanOut())
	}
}

func TestDemuxExecute(t *testing.T) {
	// init the victim
	victim := streamer.NewIndexedChannelDemux(3, MockIndexFunction)

	// prepare the scenario
	input := make(chan streamer.Message, 1)
	defer close(input)

	msg := streamer.NewMessage()
	msg.Put("testkey", "testvalue")

	input <- msg

	victim.Execute(input)
	output := <-victim.Output(1)

	if output == nil {
		t.Errorf("Expected nothing, found %s", output)
	}

	if output.Get("testkey") != "testvalue" {
		t.Errorf("Expected testvalue, found %s", output.Get("testkey"))
	}
}

func MockIndexFunction(nChannels int, input streamer.Message) int {
	return 1
}
