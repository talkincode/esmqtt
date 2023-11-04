package models

import (
	"encoding/json"
	"time"
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
	AtTime    time.Time       `json:"@timestamp,omitempty"`
	ClientID  string          `json:"clientid,omitempty"`
	Node      int64           `json:"Node,omitempty"`
}

func (e *ElasticMessage) UnmarshalJSON(data []byte) error {
	type Alias ElasticMessage
	aux := &struct {
		AtTime string `json:"@timestamp"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.AtTime != "" {
		t, err := time.Parse(time.RFC3339, aux.AtTime)
		if err != nil {
			return err
		}
		e.AtTime = t
	}
	return nil
}

func (e ElasticMessage) MarshalJSON() ([]byte, error) {
	type Alias ElasticMessage
	return json.Marshal(&struct {
		AtTime string `json:"@timestamp"`
		*Alias
	}{
		AtTime: e.AtTime.Format(time.RFC3339),
		Alias:  (*Alias)(&e),
	})
}
