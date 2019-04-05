package streamer_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/picadoh/gostreamer/streamer"
)

func TestReadLines(t *testing.T) {
	filename := "streamer_test.txt"
	contents := []byte("hello world\nhow are you?")
	err := ioutil.WriteFile(filename, contents, 0644)
	if err != nil {
		panic(err)
	}

	readContents, _ := streamer.LoadTextFile(filename)

	if readContents[0] != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'\n", readContents)
	}

	if readContents[1] != "how are you?" {
		t.Errorf("Expected 'how are you?', got '%s'\n", readContents)
	}

	os.Remove(filename)
}

func TestEmptyFile(t *testing.T) {
	filename := "streamer_test.txt"
	ioutil.WriteFile(filename, []byte{}, 0644)

	readContents, err := streamer.LoadTextFile(filename)

	if err != nil {
		t.Errorf("Expected file %s, but error occurred: %s\n", filename, err)
	}

	if len(readContents) != 0 {
		t.Errorf("Expected 0 lines, found %d\n", len(readContents))
	}

	os.Remove(filename)
}

func TestMissingFile(t *testing.T) {
	lines, err := streamer.LoadTextFile("streamer_missing.txt")

	if err == nil {
		t.Errorf("Expected file not found error, found %s\n", lines)
	}
}
