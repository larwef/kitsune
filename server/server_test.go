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

type RepositoryMock struct {
	addMessageHandler func(*kitsune.Message) error
	getMessageHandler func(string, string) (*kitsune.Message, error)
}

func (r *RepositoryMock) AddMessage(message *kitsune.Message) error {
	return r.addMessageHandler(message)
}

func (r *RepositoryMock) GetMessage(topic, id string) (*kitsune.Message, error) {
	return r.getMessageHandler(topic, id)
}

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

	server := NewServer(repo)

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/topic/topic1", bytes.NewBuffer(paylaod))
	assert.NoError(t, err)

	server.GetRouter().ServeHTTP(res, req)

	var message *kitsune.Message
	err = json.NewDecoder(res.Body).Decode(&message)
	assert.NoError(t, err)

	assert.Equal(t, message.ID, "someId")
	assert.Equal(t, message.PublishedTime, timeFromStr("2020-06-21T14:52:11.123456Z"))
	assert.Equal(t, message.Topic, "topic1")
	assert.Equal(t, message.Payload, "Some payload")
}

func TestServer_Publish_Duplicate(t *testing.T) {
	repo := &RepositoryMock{
		addMessageHandler: func(message *kitsune.Message) error {
			return repository.ErrDuplicate
		},
	}

	payload, err := json.Marshal(&PublishRequest{
		Payload: "Some payload",
	})
	assert.NoError(t, err)

	server := NewServer(repo)

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/topic/topic1", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	server.GetRouter().ServeHTTP(res, req)

	resPayload, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, res.Code)
	assert.Equal(t, "Duplicate message id\n", string(resPayload))
}

func TestServer_GetMessage(t *testing.T) {
	repo := &RepositoryMock{
		getMessageHandler: func(topic, id string) (message *kitsune.Message, err error) {
			assert.Equal(t, "topic1", topic)
			assert.Equal(t, "message1", id)
			return &kitsune.Message{}, nil
		},
	}

	server := NewServer(repo)

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/topic/topic1/message1", nil)
	assert.NoError(t, err)

	server.GetRouter().ServeHTTP(res, req)

	var message *kitsune.Message
	err = json.NewDecoder(res.Body).Decode(&message)
	assert.NoError(t, err)
}

func TestServer_GetMessage_TopicNotFound(t *testing.T) {
	repo := &RepositoryMock{
		getMessageHandler: func(topic, id string) (message *kitsune.Message, err error) {
			return nil, repository.ErrTopicNotFound
		},
	}

	server := NewServer(repo)

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/topic/topic1/message1", nil)
	assert.NoError(t, err)

	server.GetRouter().ServeHTTP(res, req)

	resPayload, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, "Topic not found\n", string(resPayload))
}

func TestServer_GetMessage_MessageNotFound(t *testing.T) {
	repo := &RepositoryMock{
		getMessageHandler: func(topic, id string) (message *kitsune.Message, err error) {
			return nil, repository.ErrMessageNotFound
		},
	}

	server := NewServer(repo)

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/topic/topic1/message1", nil)
	assert.NoError(t, err)

	server.GetRouter().ServeHTTP(res, req)

	resPayload, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, "Message not found\n", string(resPayload))
}

func timeFromStr(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic("not able to parse time from string: " + s)
	}

	return t
}
