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

type Repository struct {
	topics        map[string]*topic
	messages      map[string]*kitsune.Message
	subscriptions map[string]*subscription
}

func NewRepository() *Repository {
	return &Repository{
		topics:        map[string]*topic{},
		messages:      map[string]*kitsune.Message{},
		subscriptions: map[string]*subscription{},
	}
}

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

func (r *Repository) RetrieveMessage(topic, id string) (*kitsune.Message, error) {
	message, exists := r.messages[id]
	if !exists {
		return nil, kitsune.ErrMessageNotFound
	}

	return message, nil
}

func (r *Repository) GetMessagesFromTopic(topicName string, req kitsune.PollRequest) ([]*kitsune.Message, error) {
	s, subscriptionExists := r.subscriptions[req.SubscriptionName]
	if !subscriptionExists {
		topic, topicExists := r.topics[topicName]
		if !topicExists {
			return nil, kitsune.ErrTopicNotFound
		}

		r.subscriptions[req.SubscriptionName] = &subscription{
			topic: topic,
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

func min(a, b int) int {
	if b < a {
		return b
	}

	return a
}
