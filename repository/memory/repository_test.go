package memory

import (
	"github.com/larwef/kitsune"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRepository_PersistMessage(t *testing.T) {
	repo := Repository{
		messages: map[string]*kitsune.Message{},
		topics:   map[string]*topic{},
	}

	err := repo.PersistMessage(&kitsune.Message{
		ID:            "testId1",
		PublishedTime: time.Now(),
		Topic:         "testTopic1",
		Payload:       "testPayload1",
	})
	assert.NoError(t, err)
	assert.Len(t, repo.messages, 1)
	assert.Len(t, repo.topics, 1)
	assert.Len(t, repo.topics["testTopic1"].messages, 1)

	err = repo.PersistMessage(&kitsune.Message{
		ID:            "testId2",
		PublishedTime: time.Now(),
		Topic:         "testTopic1",
		Payload:       "testPayload2",
	})
	assert.NoError(t, err)
	assert.Len(t, repo.messages, 2)
	assert.Len(t, repo.topics, 1)
	assert.Len(t, repo.topics["testTopic1"].messages, 2)

	err = repo.PersistMessage(&kitsune.Message{
		ID:            "testId3",
		PublishedTime: time.Now(),
		Topic:         "testTopic2",
		Payload:       "testPayload3",
	})
	assert.NoError(t, err)
	assert.Len(t, repo.messages, 3)
	assert.Len(t, repo.topics, 2)
	assert.Len(t, repo.topics["testTopic1"].messages, 2)
	assert.Len(t, repo.topics["testTopic2"].messages, 1)
}

func TestRepository_RetrieveMessage(t *testing.T) {
	repo := Repository{
		messages: map[string]*kitsune.Message{
			"testId": {
				ID:            "testId",
				PublishedTime: time.Now(),
				Topic:         "testTopic",
				Payload:       "testPayload",
			}},
	}

	message, err := repo.RetrieveMessage("testTopic", "testId")
	assert.NoError(t, err)
	assert.NotNil(t, message)
}

func TestRepository_GetMessagesFromTopic(t *testing.T) {
	repo := &Repository{
		topics: map[string]*topic{
			"testTopic1": {[]*kitsune.Message{
				{ID: "testId1", PublishedTime: time.Now(), Payload: "testPayload1"},
				{ID: "testId2", PublishedTime: time.Now(), Payload: "testPayload2"},
				{ID: "testId3", PublishedTime: time.Now(), Payload: "testPayload3"},
				{ID: "testId4", PublishedTime: time.Now(), Payload: "testPayload4"},
				{ID: "testId5", PublishedTime: time.Now(), Payload: "testPayload5"},
			}},
		},
		subscriptions: map[string]*subscription{},
	}

	pollReq := kitsune.PollRequest{
		SubscriptionName:    "testSubscription",
		MaxNumberOfMessages: 2,
	}

	messages, err := repo.GetMessagesFromTopic("testTopic1", pollReq)
	assert.NoError(t, err)
	assert.Len(t, repo.subscriptions, 1)
	assert.Len(t, messages, 2)
	assert.Equal(t, "testId1", messages[0].ID)
	assert.Equal(t, "testId2", messages[1].ID)

	messages, err = repo.GetMessagesFromTopic("testTopic1", pollReq)
	assert.NoError(t, err)
	assert.Len(t, repo.subscriptions, 1)
	assert.Len(t, messages, 2)
	assert.Equal(t, "testId3", messages[0].ID)
	assert.Equal(t, "testId4", messages[1].ID)

	messages, err = repo.GetMessagesFromTopic("testTopic1", pollReq)
	assert.NoError(t, err)
	assert.Len(t, repo.subscriptions, 1)
	assert.Len(t, messages, 1)
	assert.Equal(t, "testId5", messages[0].ID)

	messages, err = repo.GetMessagesFromTopic("testTopic1", pollReq)
	assert.NoError(t, err)
	assert.Len(t, repo.subscriptions, 1)
	assert.Len(t, messages, 0)
}
