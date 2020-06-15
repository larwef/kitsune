package kitsune

import "time"

type PublishRequest struct {
	Properties map[string]string `json:"properties"`
	EventTime  *time.Time        `json:"eventTime"`
	Payload    string            `json:"payload"`
}

type Message struct {
	ID            string            `json:"id"`
	PublishedTime time.Time         `json:"publishedTime"`
	Properties    map[string]string `json:"properties,omitempty"`
	EventTime     *time.Time        `json:"eventTime,omitempty"`
	Topic         string            `json:"topic"`
	Payload       string            `json:"payload"`
}
