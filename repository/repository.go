package repository

import (
	"errors"

	"github.com/larwef/kitsune"
)

var (
	// ErrDuplicate is returned if a resource already exists.
	ErrDuplicate = errors.New("duplicate id")
	// ErrMessageNotFound is returned when message resource cant be found.
	ErrMessageNotFound = errors.New("message not found")
	// ErrTopicNotFound is returned when a topic doesnt exist.
	ErrTopicNotFound = errors.New("topic not found")
)

// MessageRepository defines the behaviour to be satified by a message repository.
type MessageRepository interface {
	AddMessage(message *kitsune.Message) error
	GetMessage(id string) (*kitsune.Message, error)
}

// TopicRepository defines the behaviour to be satified by a topic repository.
type TopicRepository interface {
	GetTopics() ([]*kitsune.Topic, error)
	GetTopic(id string) (*kitsune.Topic, error)
}

// SubscriptionRepository defines the behaviour to be satified by a subscription repository.
type SubscriptionRepository interface{}
