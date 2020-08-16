package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/larwef/kitsune"
	"github.com/larwef/kitsune/repository"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Used to simplify testing.
var now = time.Now
var id = func() string { return uuid.New().String() }

// Server holds all the http handlers.
type Server struct {
	MessageRepo      repository.MessageRepository
	TopicRepo        repository.TopicRepository
	SubscriptionRepo repository.SubscriptionRepository
}

func (s *Server) addMessage() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		var publishReq PublishRequest
		if err := receiveJSON(res, req, &publishReq); err != nil {
			zap.S().Errorw("Error marshalling request", "error", err)
			return
		}

		message := kitsune.Message{
			ID:            id(),
			PublishedTime: now(),
			Properties:    publishReq.Properties,
			EventTime:     publishReq.EventTime,
			Payload:       publishReq.Payload,
		}

		if err := s.MessageRepo.AddMessage(&message); err != nil {
			zap.S().Errorw("Error persisting message", "error", err)
			switch err {
			case repository.ErrDuplicate:
				http.Error(res, "Duplicate message id", http.StatusConflict)
			default:
				http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			}

			return
		}

		// TODO: Add message to topics.

		if err := sendJSON(res, &message); err != nil {
			zap.S().Errorw("Error marshalling response", "error", err)
			return
		}

		zap.S().Infow("Message succsessfully published",
			"id", message.ID,
			"publishedTime", message.PublishedTime,
		)
	}
}

func (s *Server) getMessage() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		params := httprouter.ParamsFromContext(req.Context())
		messageID := params.ByName("messageId")

		message, err := s.MessageRepo.GetMessage(messageID)
		if err != nil {
			zap.S().Errorw("Error retrieving message", "error", err)
			switch err {
			case repository.ErrTopicNotFound:
				http.Error(res, "Topic not found", http.StatusNotFound)
			case repository.ErrMessageNotFound:
				http.Error(res, "Message not found", http.StatusNotFound)
			default:
				http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		if err := sendJSON(res, &message); err != nil {
			zap.S().Errorw("Error marshalling response", "error", err)
			return
		}
	}
}

func (s *Server) getTopics() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		topics, err := s.TopicRepo.GetTopics()
		if err != nil {
			zap.S().Errorw("Error retrieving topics", "error", err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := sendJSON(res, &topics); err != nil {
			zap.S().Errorw("Error marshalling response", "error", err)
			return
		}
	}
}

func (s *Server) getTopic() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		params := httprouter.ParamsFromContext(req.Context())
		topicID := params.ByName("topicId")

		topic, err := s.TopicRepo.GetTopic(topicID)
		if err != nil {
			zap.S().Errorw("Error retrieving topic", "error", err)
			switch err {
			case repository.ErrTopicNotFound:
				http.Error(res, "Topic not found", http.StatusNotFound)
			default:
				http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		if err := sendJSON(res, &topic); err != nil {
			zap.S().Errorw("Error marshalling response", "error", err)
			return
		}
	}
}

func receiveJSON(res http.ResponseWriter, req *http.Request, v interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(v); err != nil {
		http.Error(res, "Unable to unmarshal request", http.StatusBadRequest)
		return err
	}

	return nil
}

func sendJSON(res http.ResponseWriter, v interface{}) error {
	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(v); err != nil {
		http.Error(res, "Error marshalling response", http.StatusInternalServerError)
		return err
	}

	return nil
}
