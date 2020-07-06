package memory

import (
	"github.com/larwef/kitsune/repository"
	"testing"
	"time"

	"github.com/larwef/kitsune"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetTopics(t *testing.T) {
	repo := Repository{
		topics: map[string]*kitsune.Topic{
			"topic1": {ID: "topic1"},
			"topic2": {ID: "topic2"},
			"topic3": {ID: "topic3"},
			"topic4": {ID: "topic4"},
			"topic5": {ID: "topic5"},
		},
	}

	topics, err := repo.GetTopics()
	assert.NoError(t, err)

	assert.Len(t, topics, 5)
}

func TestRepository_GetTopics_Empty(t *testing.T) {
	repo := Repository{
		topics: map[string]*kitsune.Topic{},
	}

	topics, err := repo.GetTopics()
	assert.NoError(t, err)

	assert.Len(t, topics, 0)
}

func TestRepository_GetTopic(t *testing.T) {
	repo := Repository{
		topics: map[string]*kitsune.Topic{
			"topic1": {ID: "topic1"},
			"topic2": {ID: "topic2"},
			"topic3": {ID: "topic3"},
			"topic4": {ID: "topic4"},
			"topic5": {ID: "topic5"},
		},
	}

	topic1, err := repo.GetTopic("topic1")
	assert.NoError(t, err)
	assert.Equal(t, topic1.ID, "topic1")

	topic3, err := repo.GetTopic("topic3")
	assert.NoError(t, err)
	assert.Equal(t, topic3.ID, "topic3")

	topic5, err := repo.GetTopic("topic5")
	assert.NoError(t, err)
	assert.Equal(t, topic5.ID, "topic5")

}

func TestRepository_GetTopic_TopicNotFound(t *testing.T) {
	repo := Repository{
		topics: map[string]*kitsune.Topic{},
	}

	_, err := repo.GetTopic("topic1")
	assert.Equal(t, repository.ErrTopicNotFound, err)
}

func TestRepository_AddMessage(t *testing.T) {
	repo := Repository{
		topics:   map[string]*kitsune.Topic{},
		messages: map[string][]*kitsune.Message{},
	}

	err := repo.AddMessage(&kitsune.Message{
		ID:            "message1",
		PublishedTime: time.Now(),
		Topic:         "topic1",
		Payload:       "payload1",
	})

	assert.NoError(t, err)
	assert.Len(t, repo.topics, 1)
	assert.Len(t, repo.messages, 1)
	assert.Len(t, repo.messages["topic1"], 1)

	err = repo.AddMessage(&kitsune.Message{
		ID:            "message2",
		PublishedTime: time.Now(),
		Topic:         "topic1",
		Payload:       "payload2",
	})
	assert.Len(t, repo.topics, 1)
	assert.Len(t, repo.messages, 1)
	assert.Len(t, repo.messages["topic1"], 2)

	err = repo.AddMessage(&kitsune.Message{
		ID:            "message3",
		PublishedTime: time.Now(),
		Topic:         "topic2",
		Payload:       "payload3",
	})
	assert.Len(t, repo.topics, 2)
	assert.Len(t, repo.messages, 2)
	assert.Len(t, repo.messages["topic1"], 2)
	assert.Len(t, repo.messages["topic2"], 1)
}

func TestRepository_AddMessage_DuplicateMessageId(t *testing.T) {
	repo := Repository{
		topics: map[string]*kitsune.Topic{
			"topic1": {ID: "topic1"},
		},
		messages: map[string][]*kitsune.Message{
			"topic1": {
				{
					ID:            "message1",
					PublishedTime: time.Now(),
					Topic:         "topic1",
					Payload:       "payload1",
				},
			},
		},
	}

	err := repo.AddMessage(&kitsune.Message{
		ID:            "message1",
		PublishedTime: time.Now(),
		Topic:         "topic1",
		Payload:       "payload1",
	})

	assert.Equal(t, repository.ErrDuplicate, err)
}

func TestRepository_GetMessage(t *testing.T) {
	repo := Repository{
		topics: map[string]*kitsune.Topic{
			"topic1": {ID: "topic1"},
			"topic2": {ID: "topic2"},
		},
		messages: map[string][]*kitsune.Message{
			"topic1": {
				{
					ID:            "message1",
					PublishedTime: time.Now(),
					Topic:         "topic1",
					Payload:       "payload1",
				},
			},
			"topic2": {
				{
					ID:            "message2",
					PublishedTime: time.Now(),
					Topic:         "topic2",
					Payload:       "payload2",
				},
			},
		},
	}

	message1, err := repo.GetMessage("topic1", "message1")
	assert.NoError(t, err)
	assert.Equal(t, "message1", message1.ID)
	assert.Equal(t, "topic1", message1.Topic)
	assert.Equal(t, "payload1", message1.Payload)

	message2, err := repo.GetMessage("topic2", "message2")
	assert.NoError(t, err)
	assert.Equal(t, "message2", message2.ID)
	assert.Equal(t, "topic2", message2.Topic)
	assert.Equal(t, "payload2", message2.Payload)
}

func TestRepository_GetMessage_TopicDoesntExist(t *testing.T) {
	repo := Repository{
		topics:   map[string]*kitsune.Topic{},
		messages: map[string][]*kitsune.Message{},
	}

	_, err := repo.GetMessage("topic1", "someId")

	assert.Equal(t, repository.ErrTopicNotFound, err)
}

func TestRepository_GetMessage_MessageDoesntExist(t *testing.T) {
	repo := Repository{
		topics: map[string]*kitsune.Topic{
			"topic1": {ID: "topic1"},
		},
		messages: map[string][]*kitsune.Message{},
	}

	_, err := repo.GetMessage("topic1", "someId")

	assert.Equal(t, repository.ErrMessageNotFound, err)
}
