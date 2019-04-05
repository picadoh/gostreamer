package streamer

import (
	"fmt"
	"strconv"
	"strings"
)

/**
Config error structure represents a configuration error message.
*/
type ConfigError struct {
	error
	message string
}

/**
Config interface provides the means to access configuration.
*/
type Config interface {
	GetString(key string) string
	GetInt(key string) int
	ToString() string
}

/**
Properties config is a key/value pair based configuration structure.
*/
type PropertiesConfig struct {
	Config
	properties map[string]string
}

/**
Loads the properties into a properties configuration instance. May return the
configuration itself along with an error that indicates if there was a problem loading the configuration.
*/
func LoadProperties(filename string) (Config, error) {
	lines, err := LoadTextFile(filename)

	if err != nil {
		return nil, nil
	}

	var raw map[string]string = make(map[string]string)

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if strings.HasPrefix(line, "#") || len(line) == 0 {
			// ignore comment and empty lines
			continue
		}

		pair := strings.Split(line, "=")

		if len(pair) != 2 {
			// invalid property format
			return nil, &ConfigError{message: fmt.Sprintf("invalid property format: %s", pair)}
		}

		key := strings.TrimSpace(pair[0])
		value := strings.TrimSpace(pair[1])

		raw[key] = value
	}

	config := &PropertiesConfig{properties: raw}
	return config, nil
}

/**
Gets a configured string value.
*/
func (config *PropertiesConfig) GetString(key string) string {
	return (*config).properties[key]
}

/**
Gets a configured int value.
*/
func (config *PropertiesConfig) GetInt(key string) int {
	value, _ := strconv.Atoi(config.properties[key])
	return value
}

/**
Dumps the configuration internal map into a string.
*/
func (config *PropertiesConfig) ToString() string {
	return fmt.Sprintf("%s", config.properties)
}

/**
Returns the error message.
*/
func (configError *ConfigError) Error() string {
	return configError.message
}

/**
Builds a new properties configuration structure
*/
func NewPropertiesConfig() Config {
	return &PropertiesConfig{properties: make(map[string]string)}
}
