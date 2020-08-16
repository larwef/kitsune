package kitsune

import (
	"time"
)

// Message represents a message resource. PublishedTime is set by server, while EventTime can be set by the publisher.
type Message struct {
	ID            string            `json:"id"`
	PublishedTime time.Time         `json:"publishedTime"`
	Properties    map[string]string `json:"properties,omitempty"`
	EventTime     *time.Time        `json:"eventTime,omitempty"`
	Payload       string            `json:"payload"`
}
