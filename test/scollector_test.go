package streamer_test

import (
	"github.com/picadoh/gostreamer/streamer"
	"testing"
)

func TestCollectData(t *testing.T) {

	done := make(chan bool)

	cfg := streamer.NewPropertiesConfig()

	victim := streamer.NewCollector("x", cfg, MockCollect)

	var output <- chan streamer.Message
	go func() {
		output = victim.Execute()
		done <- true
	}()

	<-done

	close(done)

	msg := <-output

	if msg == nil {
		t.Error("Expected output, found nothing")
	}

	if msg.Get("test") != "hello world" {
		t.Errorf("Expected 'hello world', found '%s'", msg.Get("test"))
	}
}

func MockCollect(name string, cfg streamer.Config, out *chan streamer.Message) {
	out_message := streamer.NewMessage()
	out_message.Put("test", "hello world")
	*out <- out_message
}
