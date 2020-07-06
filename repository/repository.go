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

// Repository defines the behaviour to be satified by a repository.
type Repository interface {
	TopicRepository
	SubscriptionRepository
}

// TopicRepository defines the behaviour to be satified by a repository.
type TopicRepository interface {
	GetTopics() ([]*kitsune.Topic, error)
	GetTopic(topic string) (*kitsune.Topic, error)
	AddMessage(message *kitsune.Message) error
	GetMessage(topic, id string) (*kitsune.Message, error)
}

// SubscriptionRepository defines the behaviour to be satified by a repository.
type SubscriptionRepository interface{}
