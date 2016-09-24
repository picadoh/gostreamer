package streamer_test

import (
	"testing"
	"github.com/picadoh/gostreamer/streamer"
)

func TestHash(t *testing.T) {
	hash := streamer.Hash("xpto", 5)
	if (hash != uint32(0)) {
		t.Errorf("Expected 0, found %d\n", hash)
	}

	hash = streamer.Hash("foo bar test with some hash value", 1000)
	if (hash != uint32(891)) {
		t.Errorf("Expected 0, found %d\n", hash)
	}

	hash = streamer.Hash("my hashing test", 5)
	if (hash != uint32(3)) {
		t.Errorf("Expected 3, found %d\n", hash)
	}
}
