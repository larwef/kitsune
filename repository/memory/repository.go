package memory

import (
	"github.com/larwef/kitsune"
	"github.com/larwef/kitsune/repository"
)

// Repository is a simple in memory repository.
type Repository struct {
	topics   map[string]*kitsune.Topic
	messages map[string][]*kitsune.Message
}

// NewRepository returns a new in memory Repository.
func NewRepository() *Repository {
	return &Repository{
		topics:   map[string]*kitsune.Topic{},
		messages: map[string][]*kitsune.Message{},
	}
}

// AddMessage persists a message in the repository.
func (r *Repository) AddMessage(message *kitsune.Message) error {
	if _, exists := r.topics[message.Topic]; !exists {
		r.topics[message.Topic] = &kitsune.Topic{
			ID: message.Topic,
		}
	}

	for _, m := range r.messages[message.Topic] {
		if m.ID == message.ID {
			return repository.ErrDuplicate
		}
	}

	r.messages[message.Topic] = append(r.messages[message.Topic], message)

	return nil
}

// GetMessage retrieves a spesific message from the repository.
func (r *Repository) GetMessage(topic, id string) (*kitsune.Message, error) {
	if _, exists := r.topics[topic]; !exists {
		return nil, repository.ErrTopicNotFound
	}

	for _, message := range r.messages[topic] {
		if message.ID == id {
			return message, nil
		}
	}

	return nil, repository.ErrMessageNotFound
}
