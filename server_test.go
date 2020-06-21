package kitsune

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_publish_successful(t *testing.T) {
	server := New()

	marshal, err := json.Marshal(&PublishRequest{
		Payload: "testPayload",
	})
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/publish/testTopic", bytes.NewBuffer(marshal))
	assert.NoError(t, err)
	res := httptest.NewRecorder()

	server.GetRouter().ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Len(t, server.messages, 1)
	assert.Len(t, server.topics, 1)
	assert.Len(t, server.subscriptions, 0)
}

func TestServer_getMessage(t *testing.T) {
	server := New()
	server.messages["testId"] = &Message{
		ID:            "testId",
		PublishedTime: time.Now(),
		Topic:         "testTopic",
		Payload:       "testPayload",
	}

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/testTopic/testId", nil)
	assert.NoError(t, err)

	server.GetRouter().ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestServer_poll(t *testing.T) {
	server := New()
	server.topics["testTopic"] = &Topic{
		Messages: []*Message{
			{
				ID:            "testId1",
				PublishedTime: time.Now(),
				Topic:         "testTopic",
				Payload:       "testPayload1",
			},
			{
				ID:            "testId2",
				PublishedTime: time.Now(),
				Topic:         "testTopic",
				Payload:       "testPayload2",
			},
			{
				ID:            "testId3",
				PublishedTime: time.Now(),
				Topic:         "testTopic",
				Payload:       "testPayload3",
			},
			{
				ID:            "testId4",
				PublishedTime: time.Now(),
				Topic:         "testTopic",
				Payload:       "testPayload4",
			},
			{
				ID:            "testId5",
				PublishedTime: time.Now(),
				Topic:         "testTopic",
				Payload:       "testPayload5",
			},
		},
	}

	marshal, err := json.Marshal(&PollRequest{
		SubscriptionName:    "testTopic",
		MaxNumberOfMessages: 2,
	})
	assert.NoError(t, err)

	var result []*Message

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/poll/testTopic", bytes.NewBuffer(marshal))
	assert.NoError(t, err)
	server.GetRouter().ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Len(t, server.subscriptions, 1)
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&result))
	assert.Len(t, result, 2)

	res = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodPost, "/poll/testTopic", bytes.NewBuffer(marshal))
	assert.NoError(t, err)
	server.GetRouter().ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Len(t, server.subscriptions, 1)
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&result))
	assert.Len(t, result, 2)

	res = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodPost, "/poll/testTopic", bytes.NewBuffer(marshal))
	server.GetRouter().ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Len(t, server.subscriptions, 1)
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&result))
	assert.Len(t, result, 1)

	res = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodPost, "/poll/testTopic", bytes.NewBuffer(marshal))
	server.GetRouter().ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Len(t, server.subscriptions, 1)
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&result))
	assert.Len(t, result, 0)

}
