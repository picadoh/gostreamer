package streamer_test

import (
	"testing"
	"github.com/picadoh/gostreamer/streamer"
	"math/rand"
	"time"
)

func TestRandomDemuxContext(t *testing.T) {
	mockMsg := streamer.NewMessage()

	// test with random seed
	rand.Seed(time.Now().Unix())
	index := streamer.RandomIndex(5, mockMsg)
	if index < 0 || index >= 5 {
		t.Errorf("Expected index in range [0,5[, found %d", index)
	}
}

func TestGroupDemuxContext(t *testing.T) {
	demuxCtx := streamer.NewGroupDemux("testkey")
	mockMsg := streamer.NewMessage()
	mockMsg.Put("testkey", "testvalue")

	if demuxCtx == nil {
		t.Error("Expected output, found nothing")
	}

	index := demuxCtx.GroupIndex(5, mockMsg)
	if index != 2 {
		t.Errorf("Expected 2, found %d", index)
	}

	index = demuxCtx.GroupIndex(1, mockMsg)
	if index != 0 {
		t.Errorf("Expected 0, found %d", index)
	}
}