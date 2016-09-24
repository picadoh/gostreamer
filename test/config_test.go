package streamer_test

import (
	"testing"
	"io/ioutil"
	"github.com/picadoh/gostreamer/streamer"
	"os"
)

func TestLoadProperties(t *testing.T) {
	filename := "streamer_config.cfg"
	contents := []byte("mystring=xpto\n#comment\n\nmyint=5")
	err := ioutil.WriteFile(filename, contents, 0644)
	if err != nil {
		panic(err)
	}

	config, _ := streamer.LoadProperties(filename)

	if (config.GetString("mystring") != "xpto") {
		t.Errorf("Expected 'myvalue', got '%s'\n", config.GetString("mystring"))
	}

	if (config.GetInt("myint") != 5) {
		t.Errorf("Expected 5, got '%d'\n", config.GetInt("myint"))
	}

	os.Remove(filename)
}

func TestFailedLoadProperties(t *testing.T) {
	filename := "streamer_config.cfg"
	contents := []byte("mystring=xpto\nmyint")
	err := ioutil.WriteFile(filename, contents, 0644)
	if err != nil {
		panic(err)
	}

	config, err := streamer.LoadProperties(filename)

	if (err == nil) {
		t.Errorf("Expected 'error', got config '%s'\n", config)
	}

	os.Remove(filename)
}
