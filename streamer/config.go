package streamer

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
	"io"
)

type ConfigError struct {
	error
	Message string
}

type Config interface {
	GetString(key string) string
	GetInt(key string) int
	ToString() string
}

type PropertiesConfig struct {
	Config
	Properties map[string]string
}

func LoadProperties(filename string) (Config, error) {
	var raw map[string]string = make(map[string]string)

	if len(filename) == 0 {
		return nil, nil
	}

	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if (err == io.EOF) {
			// all set
			break
		}

		if (err != nil) {
			// error, stop
			return nil, err
		}

		line = strings.TrimSpace(line)

		if (strings.HasPrefix(line, "#") || len(line) == 0) {
			// ignore comment and empty lines
			continue
		}

		pair := strings.Split(line, "=");

		if (len(pair) != 2) {
			// invalid property format
			return nil, &ConfigError{Message:fmt.Sprintf("invalid property format: %s", pair)}
		}

		key := strings.TrimSpace(pair[0])
		value := strings.TrimSpace(pair[1])

		raw[key] = value
	}

	config := &PropertiesConfig{Properties:raw}
	return config, nil
}

func (config *PropertiesConfig) GetString(key string) string {
	return (*config).Properties[key]
}

func (config *PropertiesConfig) GetInt(key string) int {
	value, _ := strconv.Atoi(config.Properties[key])
	return value
}

func (config *PropertiesConfig) ToString() string {
	return fmt.Sprintf("%s", config.Properties)
}

func (configError *ConfigError) Error() string {
	return configError.Message
}