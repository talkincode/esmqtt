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
	ID        string          `json:"id,omitempty"`
	Index     string          `json:"index,omitempty"`
	Vector    []float32       `json:"vector,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp int64           `json:"timestamp,omitempty"`
}
