package streamer_test

import (
	"github.com/picadoh/gostreamer/streamer"
	"testing"
	"log"
)

func TestCollectData(t *testing.T) {
	cfg := streamer.NewPropertiesConfig()

	victim := streamer.NewCollector("x", cfg, MockCollect)

	output := victim.Execute()

	msg := <-output

	if msg == nil {
		t.Error("Expected output, found nothing")
	}

	if msg.Get("test") != "hello world" {
		t.Errorf("Expected 'hello world', found '%s'", msg.Get("test"))
	}
}

func MockCollect(name string, cfg streamer.Config, out chan streamer.Message) {
	log.Printf("Executing mocked collect %s with config %s", name, cfg.ToString())

	out_message := streamer.NewMessage()
	out_message.Put("test", "hello world")
	out <- out_message
}
