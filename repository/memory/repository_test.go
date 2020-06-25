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

func TestRepository_PersistMessage_Duplicate(t *testing.T) {
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

	err = repo.PersistMessage(&kitsune.Message{
		ID:            "testId1",
		PublishedTime: time.Now(),
		Topic:         "testTopic1",
		Payload:       "testPayload1",
	})
	assert.Equal(t, kitsune.ErrDuplicateMessage, err)
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

	message, err := repo.GetMessage("testTopic", "testId")
	assert.NoError(t, err)
	assert.NotNil(t, message)
}

func TestRepository_RetrieveMessage_MessageDoesntExist(t *testing.T) {
	repo := Repository{
		messages: map[string]*kitsune.Message{},
	}

	_, err := repo.GetMessage("testTopic", "testId")
	assert.Equal(t, kitsune.ErrMessageNotFound, err)
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

	messages, err := repo.PollTopic("testTopic1", pollReq)
	assert.NoError(t, err)
	assert.Len(t, repo.subscriptions, 1)
	assert.Len(t, messages, 2)
	assert.Equal(t, "testId1", messages[0].ID)
	assert.Equal(t, "testId2", messages[1].ID)

	messages, err = repo.PollTopic("testTopic1", pollReq)
	assert.NoError(t, err)
	assert.Len(t, repo.subscriptions, 1)
	assert.Len(t, messages, 2)
	assert.Equal(t, "testId3", messages[0].ID)
	assert.Equal(t, "testId4", messages[1].ID)

	messages, err = repo.PollTopic("testTopic1", pollReq)
	assert.NoError(t, err)
	assert.Len(t, repo.subscriptions, 1)
	assert.Len(t, messages, 1)
	assert.Equal(t, "testId5", messages[0].ID)

	messages, err = repo.PollTopic("testTopic1", pollReq)
	assert.NoError(t, err)
	assert.Len(t, repo.subscriptions, 1)
	assert.Len(t, messages, 0)
}

func TestRepository_GetMessagesFromTopic_TopicDoesntExist(t *testing.T) {
	repo := &Repository{
		topics:        map[string]*topic{},
		subscriptions: map[string]*subscription{},
	}

	pollReq := kitsune.PollRequest{
		SubscriptionName:    "testSubscription",
		MaxNumberOfMessages: 2,
	}

	_, err := repo.PollTopic("testTopic", pollReq)
	assert.Equal(t, kitsune.ErrTopicNotFound, err)
}

func TestRepository_SetSubscriptionPosition(t *testing.T) {
	repo := &Repository{
		topics: map[string]*topic{
			"testTopic": {[]*kitsune.Message{
				{ID: "testId1", PublishedTime: timeFromStr("2020-06-24T19:00:00.000000Z"), Payload: "testPayload1"},
				{ID: "testId2", PublishedTime: timeFromStr("2020-06-24T20:00:00.000000Z"), Payload: "testPayload2"},
				{ID: "testId3", PublishedTime: timeFromStr("2020-06-24T21:00:00.000000Z"), Payload: "testPayload3"},
				{ID: "testId4", PublishedTime: timeFromStr("2020-06-25T19:00:00.000000Z"), Payload: "testPayload4"},
				{ID: "testId5", PublishedTime: timeFromStr("2020-06-25T20:00:00.000000Z"), Payload: "testPayload5"},
			}},
		},
		subscriptions: map[string]*subscription{},
	}

	posReq := kitsune.SubscriptionPositionRequest{
		SubscriptionName: "testSubscription",
		MessageID:        "someId",
	}

	// Setting a non existing subscription should create it
	err := repo.SetSubscriptionPosition("testTopic", posReq)
	assert.NoError(t, err)
	assert.Len(t, repo.subscriptions, 1)

	pollReq := kitsune.PollRequest{
		SubscriptionName:    "testSubscription",
		MaxNumberOfMessages: 3,
	}

	// Poll some messages
	messages, err := repo.PollTopic("testTopic", pollReq)
	assert.NoError(t, err)
	assert.Len(t, repo.subscriptions, 1)
	assert.Len(t, messages, 3)

	// Should set to first message since neither time or id is found
	err = repo.SetSubscriptionPosition("testTopic", posReq)
	assert.NoError(t, err)

	// Poll some messages. First message should now be the first returned
	messages, err = repo.PollTopic("testTopic", pollReq)
	assert.NoError(t, err)
	assert.Len(t, repo.subscriptions, 1)
	assert.Len(t, messages, 3)
	assert.Equal(t, "testId1", messages[0].ID)

	posReq.MessageID = "testId2"
	// Should set to second message based on id
	err = repo.SetSubscriptionPosition("testTopic", posReq)
	assert.NoError(t, err)

	// Poll some messages. Second message should now be the first returned
	messages, err = repo.PollTopic("testTopic", pollReq)
	assert.NoError(t, err)
	assert.Len(t, repo.subscriptions, 1)
	assert.Len(t, messages, 3)
	assert.Equal(t, "testId2", messages[0].ID)

	posReq.MessageID = "testId2"
	pubTime := timeFromStr("2020-06-25T19:00:00.000000Z")
	posReq.PublishedTime = &pubTime
	// Should set to fourth message based on time.
	err = repo.SetSubscriptionPosition("testTopic", posReq)
	assert.NoError(t, err)

	// Poll some messages. Fourth message should now be the first returned
	messages, err = repo.PollTopic("testTopic", pollReq)
	assert.NoError(t, err)
	assert.Len(t, repo.subscriptions, 1)
	assert.Len(t, messages, 2)
	assert.Equal(t, "testId4", messages[0].ID)
}

func TestRepository_SetSubscriptionPosition_TopicDoesntExist(t *testing.T) {
	repo := &Repository{}

	req := kitsune.SubscriptionPositionRequest{
		SubscriptionName: "testSubscription",
		MessageID:        "someId",
	}

	err := repo.SetSubscriptionPosition("testTopic", req)
	assert.Equal(t, kitsune.ErrTopicNotFound, err)
}

func timeFromStr(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic("not able to parse time from string: " + s)
	}

	return t
}
