package repository

import (
	"errors"
	"github.com/larwef/kitsune"
)

var (
	// ErrDuplicateMessage is returned if message resource already exists.
	ErrDuplicateMessage = errors.New("duplicate message id")
	// ErrTopicNotFound is returned when a topic doesnt exist.
	ErrTopicNotFound = errors.New("topic not found")
	// ErrMessageNotFound is returned when message resource cant be found.
	ErrMessageNotFound = errors.New("message not found")
	// ErrCantSpecifyIDAndTime is returned when trying to specify both time and id when setting a subscription.
	ErrCantSpecifyIDAndTime = errors.New("cant specify both a time and an id when setting a subscription position")
)

// Repository defines the behaviour to be satified by a repository.
type Repository interface {
	PersistMessage(message *kitsune.Message) error
	PollTopic(topicName string, req kitsune.PollRequest) ([]*kitsune.Message, error)
	GetMessage(topic, id string) (*kitsune.Message, error)
	SetSubscriptionPosition(topicName string, req kitsune.SubscriptionPositionRequest) error
}
