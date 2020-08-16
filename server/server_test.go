package server

import (
	"bytes"
	"encoding/json"
	"github.com/larwef/kitsune"
	"github.com/larwef/kitsune/repository"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Mock
type RepositoryMock struct {
	getTopicsHandler  func() ([]*kitsune.Topic, error)
	getTopicHandler   func(string) (*kitsune.Topic, error)
	addMessageHandler func(*kitsune.Message) error
	getMessageHandler func(string) (*kitsune.Message, error)
}

func (r *RepositoryMock) GetTopics() ([]*kitsune.Topic, error) {
	return r.getTopicsHandler()
}

func (r *RepositoryMock) GetTopic(topic string) (*kitsune.Topic, error) {
	return r.getTopicHandler(topic)
}

func (r *RepositoryMock) AddMessage(message *kitsune.Message) error {
	return r.addMessageHandler(message)
}

func (r *RepositoryMock) GetMessage(id string) (*kitsune.Message, error) {
	return r.getMessageHandler(id)
}

// Tests
func TestServer_Publish(t *testing.T) {
	now = func() time.Time {
		n, err := time.Parse(time.RFC3339, "2020-06-21T14:52:11.123456Z")
		assert.NoError(t, err)
		return n
	}

	id = func() string {
		return "someId"
	}

	repo := &RepositoryMock{
		addMessageHandler: func(message *kitsune.Message) error {
			return nil
		},
	}

	paylaod, err := json.Marshal(&PublishRequest{
		Payload: "Some payload",
	})
	assert.NoError(t, err)

	server := &Server{
		MessageRepo: repo,
	}

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/message", bytes.NewBuffer(paylaod))
	assert.NoError(t, err)

	server.GetRouter().ServeHTTP(res, req)

	var message *kitsune.Message
	err = json.NewDecoder(res.Body).Decode(&message)
	assert.NoError(t, err)

	assert.Equal(t, message.ID, "someId")
	assert.Equal(t, message.PublishedTime, timeFromStr("2020-06-21T14:52:11.123456Z"))
	assert.Equal(t, message.Payload, "Some payload")
}

func TestServer_GetMessage(t *testing.T) {
	repo := &RepositoryMock{
		getMessageHandler: func(id string) (message *kitsune.Message, err error) {
			assert.Equal(t, "message1", id)
			return &kitsune.Message{}, nil
		},
	}

	server := &Server{
		MessageRepo: repo,
	}

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/message/message1", nil)
	assert.NoError(t, err)

	server.GetRouter().ServeHTTP(res, req)

	var message *kitsune.Message
	err = json.NewDecoder(res.Body).Decode(&message)
	assert.NoError(t, err)
}

func TestServer_GetMessage_MessageNotFound(t *testing.T) {
	repo := &RepositoryMock{
		getMessageHandler: func(id string) (message *kitsune.Message, err error) {
			return nil, repository.ErrMessageNotFound
		},
	}

	server := &Server{
		MessageRepo: repo,
	}

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/message/message1", nil)
	assert.NoError(t, err)

	server.GetRouter().ServeHTTP(res, req)

	resPayload, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, "Message not found\n", string(resPayload))
}

func TestServer_GetTopics(t *testing.T) {
	repo := &RepositoryMock{
		getTopicsHandler: func() ([]*kitsune.Topic, error) {
			return []*kitsune.Topic{
				{ID: "topic1"},
				{ID: "topic2"},
				{ID: "topic3"},
				{ID: "topic4"},
				{ID: "topic5"},
			}, nil
		},
	}

	server := &Server{
		TopicRepo: repo,
	}

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/topic", nil)
	assert.NoError(t, err)

	server.GetRouter().ServeHTTP(res, req)

	var topics []*kitsune.Topic
	err = json.NewDecoder(res.Body).Decode(&topics)
	assert.NoError(t, err)
	assert.Len(t, topics, 5)
}

func TestServer_GetTopics_Empty(t *testing.T) {
	repo := &RepositoryMock{
		getTopicsHandler: func() ([]*kitsune.Topic, error) {
			return []*kitsune.Topic{}, nil
		},
	}

	server := &Server{
		TopicRepo: repo,
	}

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/topic", nil)
	assert.NoError(t, err)

	server.GetRouter().ServeHTTP(res, req)

	var topics []*kitsune.Topic
	err = json.NewDecoder(res.Body).Decode(&topics)
	assert.NoError(t, err)
	assert.Len(t, topics, 0)
}

func TestServer_GetTopic(t *testing.T) {
	repo := &RepositoryMock{
		getTopicHandler: func(topic string) (*kitsune.Topic, error) {
			assert.Equal(t, "topic1", topic)
			return &kitsune.Topic{ID: "topic1"}, nil
		},
	}

	server := &Server{
		TopicRepo: repo,
	}

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/topic/topic1", nil)
	assert.NoError(t, err)

	server.GetRouter().ServeHTTP(res, req)

	var topic *kitsune.Topic
	err = json.NewDecoder(res.Body).Decode(&topic)
	assert.NoError(t, err)
	assert.Equal(t, topic.ID, "topic1")
}

func TestServer_GetTopic_TopicNotFound(t *testing.T) {
	repo := &RepositoryMock{
		getTopicHandler: func(topic string) (*kitsune.Topic, error) {
			return nil, repository.ErrTopicNotFound
		},
	}

	server := &Server{
		TopicRepo: repo,
	}

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/topic/topic1", nil)
	assert.NoError(t, err)

	server.GetRouter().ServeHTTP(res, req)

	resPayload, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, "Topic not found\n", string(resPayload))
}

func timeFromStr(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic("not able to parse time from string: " + s)
	}

	return t
}
