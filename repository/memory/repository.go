package memory

import (
	"github.com/larwef/kitsune"
	"github.com/larwef/kitsune/repository"
)

// Repository is a simple in memory repository.
type Repository struct {
	topics   map[string]*kitsune.Topic
	messages map[string]*kitsune.Message
}

// NewRepository returns a new in memory Repository.
func NewRepository() *Repository {
	return &Repository{
		topics:   map[string]*kitsune.Topic{},
		messages: map[string]*kitsune.Message{},
	}
}

// AddMessage persists a message in the repository.
func (r *Repository) AddMessage(message *kitsune.Message) error {
	if _, exists := r.messages[message.ID]; exists {
		return repository.ErrDuplicate
	}

	r.messages[message.ID] = message

	return nil
}

// GetMessage retrieves a spesific message from the repository.
func (r *Repository) GetMessage(messageID string) (*kitsune.Message, error) {
	message, exists := r.messages[messageID]
	if !exists {
		return nil, repository.ErrMessageNotFound
	}

	return message, nil
}

// GetTopics lists all registered topics.
func (r *Repository) GetTopics() ([]*kitsune.Topic, error) {
	topics := make([]*kitsune.Topic, 0)

	for _, v := range r.topics {
		topics = append(topics, v)
	}

	return topics, nil
}

// GetTopic shows a specific topic.
func (r *Repository) GetTopic(id string) (*kitsune.Topic, error) {
	topic, exists := r.topics[id]
	if !exists {
		return nil, repository.ErrTopicNotFound
	}

	return topic, nil
}
