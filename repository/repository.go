package repository

import (
	"github.com/larwef/kitsune"
)

// Repository defines the behaviour to be satified by a repository.
type Repository interface {
	PersistMessage(message *kitsune.Message) error
	PollTopic(topicName string, req kitsune.PollRequest) ([]*kitsune.Message, error)
	GetMessage(topic, id string) (*kitsune.Message, error)
	SetSubscriptionPosition(topicName string, req kitsune.SubscriptionPositionRequest) error
}
