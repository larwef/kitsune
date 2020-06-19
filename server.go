package kitsune

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	topics   map[string]*Topic
	messages map[string]*Message
}

func New() *Server {
	return &Server{
		topics:   make(map[string]*Topic),
		messages: make(map[string]*Message),
	}
}

func (s *Server) publish() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		var publishReq PublishRequest
		if err := json.NewDecoder(req.Body).Decode(&publishReq); err != nil {
			zap.S().Errorw("Error marshalling request", "error", err)
			http.Error(res, "Unable to unmarshal request", http.StatusBadRequest)
			return
		}

		params := httprouter.ParamsFromContext(req.Context())
		message := Message{
			ID:            uuid.New().String(),
			PublishedTime: time.Now(),
			Properties:    publishReq.Properties,
			EventTime:     publishReq.EventTime,
			Topic:         params.ByName("topic"),
			Payload:       publishReq.Payload,
		}
		if err := s.persistMessage(message); err != nil {
			zap.S().Errorw("Error persisting message", "error", err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(&message); err != nil {
			zap.S().Errorw("Error marshalling response", "error", err)
			http.Error(res, "Error marshalling response", http.StatusInternalServerError)
			return
		}

		zap.S().Infow("Message succsessfully published message",
			"topic", message.Topic,
			"id", message.ID,
			"publishedTime", message.PublishedTime,
		)
	}
}

func (s *Server) getMessage() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		params := httprouter.ParamsFromContext(req.Context())

		message, err := s.retrieveMessage(params.ByName("topic"), params.ByName("messageId"))
		if err != nil {
			zap.S().Errorw("Error retrieving message", "error", err)
			// TODO: When doing this properly, make a custom error type and make sure its actually not found and not some other error
			http.Error(res, "Not Found", http.StatusNotFound)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(&message); err != nil {
			zap.S().Errorw("Error marshalling response", "error", err)
			http.Error(res, "Error marshalling response", http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) persistMessage(message Message) error {
	s.messages[message.ID] = &message

	topic, exists := s.topics[message.Topic]
	if !exists {
		s.topics[message.Topic] = &Topic{Messages: []*Message{&message}}
		return nil
	}

	topic.Messages = append(topic.Messages, &message)

	return nil
}

func (s *Server) retrieveMessage(topic, id string) (*Message, error) {
	message, exists := s.messages[id]
	if !exists {
		return nil, fmt.Errorf("could not find message %q in topic %q", id, topic)
	}

	return message, nil
}
