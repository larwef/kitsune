package kitsune

import "time"

type Event struct {
	ID            string
	PublishedTime time.Time
	Properties    map[string]string
	EventTime     time.Time
	Payload       []byte
}
