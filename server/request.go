package server

import "time"

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
