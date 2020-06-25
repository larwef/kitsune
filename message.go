package kitsune

import (
	"errors"
	"time"
)

var (
	// ErrDuplicateMessage is returned if message resource already exists.
	ErrDuplicateMessage = errors.New("duplicate message id")
	// ErrTopicNotFound is returned when a topic doesnt exist.
	ErrTopicNotFound = errors.New("topic not found")
	// ErrMessageNotFound is returned when message resource cant be found.
	ErrMessageNotFound = errors.New("message not found")
)

// PublishRequest is used to publish a message.
type PublishRequest struct {
	Properties map[string]string `json:"properties"`
	EventTime  *time.Time        `json:"eventTime"`
	Payload    string            `json:"payload"`
}

// PollRequest is used to poll messages from a topic.
type PollRequest struct {
	SubscriptionName    string `json:"subscriptionName"`
	MaxNumberOfMessages uint   `json:"maxNumberOfMessages"`
}

// SubscriptionPositionRequest is used to set . Will stop at first match on either MessageID or publishtime. Should specify the
// first message client should receive when polling.
type SubscriptionPositionRequest struct {
	SubscriptionName string     `json:"subscriptionName"`
	MessageID        string     `json:"messageId,omitempty"`
	PublishedTime    *time.Time `json:"publishedTime,omitempty"`
}

// Message represents a message resource. PublishedTime is set by server, while EventTime can be set by the publisher.
type Message struct {
	ID            string            `json:"id"`
	PublishedTime time.Time         `json:"publishedTime"`
	Properties    map[string]string `json:"properties,omitempty"`
	EventTime     *time.Time        `json:"eventTime,omitempty"`
	Topic         string            `json:"topic"`
	Payload       string            `json:"payload"`
}
