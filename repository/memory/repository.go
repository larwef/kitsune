package memory

import (
	"github.com/larwef/kitsune"
)

type topic struct {
	messages []*kitsune.Message
}

type subscription struct {
	topic *topic
	index uint
}

// Repository is a simple in memory repository.
type Repository struct {
	topics        map[string]*topic
	messages      map[string]*kitsune.Message
	subscriptions map[string]*subscription
}

// NewRepository returns a new in memory Repository.
func NewRepository() *Repository {
	return &Repository{
		topics:        map[string]*topic{},
		messages:      map[string]*kitsune.Message{},
		subscriptions: map[string]*subscription{},
	}
}

// PersistMessage persists a message in the repository.
func (r *Repository) PersistMessage(message *kitsune.Message) error {
	if _, exists := r.messages[message.ID]; exists {
		return kitsune.ErrDuplicateMessage
	}

	r.messages[message.ID] = message

	t, exists := r.topics[message.Topic]
	if !exists {
		r.topics[message.Topic] = &topic{messages: []*kitsune.Message{message}}
		return nil
	}

	t.messages = append(t.messages, message)

	return nil
}

// GetMessage retrieves a spesific message from the repository.
func (r *Repository) GetMessage(topic, id string) (*kitsune.Message, error) {
	message, exists := r.messages[id]
	if !exists {
		return nil, kitsune.ErrMessageNotFound
	}

	return message, nil
}

// PollTopic polls messages from a topic as specified in the Pollrequest.
func (r *Repository) PollTopic(topicName string, req kitsune.PollRequest) ([]*kitsune.Message, error) {
	t, topicExists := r.topics[topicName]
	if !topicExists {
		return nil, kitsune.ErrTopicNotFound
	}

	s, subscriptionExists := r.subscriptions[req.SubscriptionName]
	if !subscriptionExists {
		r.subscriptions[req.SubscriptionName] = &subscription{
			topic: t,
			index: 0,
		}

		s = r.subscriptions[req.SubscriptionName]
	}

	start := s.index
	end := min(int(s.index+req.MaxNumberOfMessages), len(s.topic.messages))

	messages := s.topic.messages[start:end]
	s.index = uint(end)

	return messages, nil
}

// SetSubscriptionPosition is used to set the subscription position to a desired message in the stream.
func (r *Repository) SetSubscriptionPosition(topicName string, req kitsune.SubscriptionPositionRequest) error {
	t, topicExists := r.topics[topicName]
	if !topicExists {
		return kitsune.ErrTopicNotFound
	}

	s, subscriptionExists := r.subscriptions[req.SubscriptionName]
	if !subscriptionExists {
		r.subscriptions[req.SubscriptionName] = &subscription{
			topic: t,
			index: 0,
		}

		s = r.subscriptions[req.SubscriptionName]
	}

	for i := range t.messages {
		index := len(t.messages) - i - 1
		if t.messages[index].ID == req.MessageID || (req.PublishedTime != nil && t.messages[index-1].PublishedTime.Before(*req.PublishedTime)) {
			s.index = uint(index)
			return nil
		}
	}

	s.index = 0
	return nil
}

func min(a, b int) int {
	if b < a {
		return b
	}

	return a
}
