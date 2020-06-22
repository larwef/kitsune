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

type Server struct {
	repo repository.Repository
}

func NewServer(repo repository.Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) publish() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		var publishReq kitsune.PublishRequest
		if err := json.NewDecoder(req.Body).Decode(&publishReq); err != nil {
			zap.S().Errorw("Error marshalling request", "error", err)
			http.Error(res, "Unable to unmarshal request", http.StatusBadRequest)
			return
		}

		params := httprouter.ParamsFromContext(req.Context())
		message := kitsune.Message{
			ID:            id(),
			PublishedTime: now(),
			Properties:    publishReq.Properties,
			EventTime:     publishReq.EventTime,
			Topic:         params.ByName("topic"),
			Payload:       publishReq.Payload,
		}

		if err := s.repo.PersistMessage(&message); err != nil {
			zap.S().Errorw("Error persisting message", "error", err)
			switch err {
			case kitsune.ErrDuplicateMessage:
				http.Error(res, "Duplicate message id", http.StatusConflict)
			default:
				http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			}

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

		message, err := s.repo.RetrieveMessage(params.ByName("topic"), params.ByName("messageId"))
		if err != nil {
			zap.S().Errorw("Error retrieving message", "error", err)
			switch err {
			case kitsune.ErrMessageNotFound:
				http.Error(res, "Message not found", http.StatusNotFound)
			default:
				http.Error(res, "Not Found", http.StatusNotFound)
			}
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

func (s *Server) poll() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		topic := httprouter.ParamsFromContext(req.Context()).ByName("topic")

		var pollReq kitsune.PollRequest
		if err := json.NewDecoder(req.Body).Decode(&pollReq); err != nil {
			zap.S().Errorw("Error marshalling request", "error", err)
			http.Error(res, "Unable to unmarshal request", http.StatusBadRequest)
			return
		}

		messages, err := s.repo.GetMessagesFromTopic(topic, pollReq)
		if err != nil {
			zap.S().Errorw("Error polling messages", "error", err)
			switch err {
			case kitsune.ErrTopicNotFound:
				http.Error(res, "Topic not found", http.StatusNotFound)
			default:
				http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			}

			return
		}

		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(&messages); err != nil {
			zap.S().Errorw("Error marshalling response", "error", err)
			http.Error(res, "Error marshalling response", http.StatusInternalServerError)
			return
		}

		zap.S().Infow("Messages sucessfully polled",
			"topic", topic,
			"subscription", pollReq.SubscriptionName,
			"messagesReturned", len(messages),
		)
	}
}
