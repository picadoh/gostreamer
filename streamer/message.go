package streamer

import "time"

/**
Message interface provides the means to get and set the values for a specific key in a key/value set.
 */
type Message interface {
	Get(key string) interface{}
	Put(key string, value interface{})
}

/**
The base message implementation is a set of key/value pairs enriched with a timestamp related to the message creation.
 */
type MessageImpl struct {
	timestamp string
	payload   map[string]interface{}
}

/**
The get method retrieves a value from a message based on its key.
 */
func (m *MessageImpl) Get(key string) interface{} {
	return m.payload[key]
}

/**
The set method defines a value for a specific key inside the message.
 */
func (m *MessageImpl) Put(key string, value interface{}) {
	m.payload[key] = value
}

/**
Builds a new message and sets the timestamp to be the current time
 */
func NewMessage() Message {
	return &MessageImpl{timestamp:time.Now().String(), payload:make(map[string]interface{})}
}
