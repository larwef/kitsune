package server

import (
	"bytes"
	"encoding/json"
	"github.com/larwef/kitsune"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type RepositoryMock struct {
	mock.Mock
	persistMessageHandler      func(*kitsune.Message) error
	retrieveMessageHandler     func(string, string) (*kitsune.Message, error)
	getMessageFromTopicHandler func(string, kitsune.PollRequest) ([]*kitsune.Message, error)
}

func (r *RepositoryMock) PersistMessage(message *kitsune.Message) error {
	return r.persistMessageHandler(message)
}

func (r *RepositoryMock) RetrieveMessage(topic, id string) (*kitsune.Message, error) {
	return r.retrieveMessageHandler(topic, id)
}

func (r *RepositoryMock) GetMessagesFromTopic(topicName string, req kitsune.PollRequest) ([]*kitsune.Message, error) {
	return r.getMessageFromTopicHandler(topicName, req)
}

func TestServer(t *testing.T) {
	tests := []struct {
		payload                    interface{}
		url                        string
		method                     string
		persistMessageHandler      func(*kitsune.Message) error
		retrieveMessageHandler     func(string, string) (*kitsune.Message, error)
		getMessageFromTopicHandler func(string, kitsune.PollRequest) ([]*kitsune.Message, error)
		expectedStatus             int
		expectedPayload            string
	}{
		// Publish
		{
			payload: kitsune.PublishRequest{
				Payload: "testPayload",
			},
			url:    "/publish/testTopic",
			method: http.MethodPost,
			persistMessageHandler: func(message *kitsune.Message) error {
				return nil
			},
			expectedStatus:  http.StatusOK,
			expectedPayload: "{\"id\":\"someId\",\"publishedTime\":\"2020-06-21T14:52:11.123456Z\",\"topic\":\"testTopic\",\"payload\":\"testPayload\"}\n",
		},
		{
			payload: kitsune.PublishRequest{
				Payload: "testPayload",
			},
			url:    "/publish/testTopic",
			method: http.MethodPost,
			persistMessageHandler: func(message *kitsune.Message) error {
				return kitsune.ErrDuplicateMessage
			},
			expectedStatus:  http.StatusConflict,
			expectedPayload: "Duplicate message id\n",
		},
		// GetMessage
		{
			payload: nil,
			url:     "/testTopic/someId",
			method:  http.MethodGet,
			retrieveMessageHandler: func(topic, id string) (*kitsune.Message, error) {
				return &kitsune.Message{
					ID:            "someId",
					PublishedTime: now(),
					Topic:         "testTopic",
					Payload:       "testPayload",
				}, nil
			},
			expectedStatus:  http.StatusOK,
			expectedPayload: "{\"id\":\"someId\",\"publishedTime\":\"2020-06-21T14:52:11.123456Z\",\"topic\":\"testTopic\",\"payload\":\"testPayload\"}\n",
		},
		{
			payload:         nil,
			url:             "/testTopic/someId",
			method:          http.MethodPost,
			expectedStatus:  http.StatusMethodNotAllowed,
			expectedPayload: "Method Not Allowed\n",
		},
		{
			payload: nil,
			url:     "/testTopic/someId",
			method:  http.MethodGet,
			retrieveMessageHandler: func(topic, id string) (*kitsune.Message, error) {
				return nil, kitsune.ErrMessageNotFound
			},
			expectedStatus:  http.StatusNotFound,
			expectedPayload: "Message not found\n",
		},
		// Poll
		{
			payload: &kitsune.PollRequest{
				SubscriptionName:    "testTopic",
				MaxNumberOfMessages: 10,
			},
			url:    "/poll/testTopic",
			method: http.MethodPost,
			getMessageFromTopicHandler: func(topic string, req kitsune.PollRequest) ([]*kitsune.Message, error) {
				return []*kitsune.Message{}, nil
			},
			expectedStatus:  http.StatusOK,
			expectedPayload: "[]\n",
		},
		{
			payload: &kitsune.PollRequest{
				SubscriptionName:    "testTopic",
				MaxNumberOfMessages: 10,
			},
			url:    "/poll/testTopic",
			method: http.MethodPost,
			getMessageFromTopicHandler: func(topic string, req kitsune.PollRequest) ([]*kitsune.Message, error) {
				return []*kitsune.Message{}, kitsune.ErrTopicNotFound
			},
			expectedStatus:  http.StatusNotFound,
			expectedPayload: "Topic not found\n",
		},
	}

	now = func() time.Time {
		n, err := time.Parse(time.RFC3339, "2020-06-21T14:52:11.123456Z")
		assert.NoError(t, err)
		return n
	}

	id = func() string {
		return "someId"
	}

	for _, test := range tests {
		marshal, err := json.Marshal(&test.payload)
		assert.NoError(t, err)

		req, err := http.NewRequest(test.method, test.url, bytes.NewBuffer(marshal))
		assert.NoError(t, err)

		res := httptest.NewRecorder()
		server := NewServer(&RepositoryMock{
			persistMessageHandler:      test.persistMessageHandler,
			retrieveMessageHandler:     test.retrieveMessageHandler,
			getMessageFromTopicHandler: test.getMessageFromTopicHandler,
		})

		server.GetRouter().ServeHTTP(res, req)

		payload, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		assert.Equal(t, test.expectedStatus, res.Code)
		assert.Equal(t, test.expectedPayload, string(payload))
	}
}
