package kitsune

import (
	"bytes"
	"encoding/json"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	tests := []struct {
		payload        interface{}
		url            string
		method         string
		expectedStatus int
	}{
		{
			payload:        nil,
			url:            "",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
		},
		{
			payload:        nil,
			url:            "/testTopic/nonExistingId",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
		},
		{
			payload:        PublishRequest{Payload: "SomePayload"},
			url:            "/notPublish",
			method:         http.MethodPost,
			expectedStatus: http.StatusNotFound,
		},
		{
			payload:        PublishRequest{Payload: "SomePayload"},
			url:            "/publish/testTopic",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			payload:        nil,
			url:            "/publish/testTopic",
			method:         http.MethodPut,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, test := range tests {
		marshal, err := json.Marshal(&test.payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(test.method, test.url, bytes.NewBuffer(marshal))
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()
		server := New()
		server.GetRouter().ServeHTTP(res, req)

		assert.Equal(t, res.Code, test.expectedStatus)
	}
}
