package repository

import (
	"github.com/larwef/kitsune"
)

type Repository interface {
	PersistMessage(message *kitsune.Message) error
	RetrieveMessage(topic, id string) (*kitsune.Message, error)
	GetMessagesFromTopic(topicName string, req kitsune.PollRequest) ([]*kitsune.Message, error)
}
