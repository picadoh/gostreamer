package main

import "time"

type Message interface {
	Get(key string) interface{}
	Put(key string, value interface{})
}

type MessageImpl struct {
	timestamp string
	payload   map[string]interface{}
}

func (m *MessageImpl) Get(key string) interface{} {
	return m.payload[key]
}

func (m *MessageImpl) Put(key string, value interface{}) {
	m.payload[key] = value
}

func NewMessage() Message {
	return &MessageImpl{timestamp:time.Now().String(), payload:make(map[string]interface{})}
}
