package models

import (
	"encoding/json"
)

type MessageRule struct {
	Topic   string `json:"topic"`
	Index   string `json:"index"`
	Spliter string `json:"spliter"`
}

type ElasticMessage struct {
	ID        string          `json:"id"`
	Index     string          `json:"index"`
	Vector    []float32       `json:"vector,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp int64           `json:"timestamp"`
}
