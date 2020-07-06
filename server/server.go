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
	repo repository.Repository
}

// NewServer returns a new Server object.
func NewServer(repo repository.Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) getTopics() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		topics, err := s.repo.GetTopics()
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

		topic, err := s.repo.GetTopic(topicID)
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

func (s *Server) addMessage() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		var publishReq PublishRequest
		if err := receiveJSON(res, req, &publishReq); err != nil {
			zap.S().Errorw("Error marshalling request", "error", err)
			return
		}

		params := httprouter.ParamsFromContext(req.Context())
		topic := params.ByName("topicId")

		message := kitsune.Message{
			ID:            id(),
			PublishedTime: now(),
			Properties:    publishReq.Properties,
			EventTime:     publishReq.EventTime,
			Topic:         topic,
			Payload:       publishReq.Payload,
		}

		if err := s.repo.AddMessage(&message); err != nil {
			zap.S().Errorw("Error persisting message", "error", err)
			switch err {
			case repository.ErrDuplicate:
				http.Error(res, "Duplicate message id", http.StatusConflict)
			default:
				http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			}

			return
		}

		if err := sendJSON(res, &message); err != nil {
			zap.S().Errorw("Error marshalling response", "error", err)
			return
		}

		zap.S().Infow("Message succsessfully published",
			"topic", message.Topic,
			"id", message.ID,
			"publishedTime", message.PublishedTime,
		)
	}
}

func (s *Server) getMessage() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		params := httprouter.ParamsFromContext(req.Context())
		topicID := params.ByName("topicId")
		messageID := params.ByName("messageId")

		message, err := s.repo.GetMessage(topicID, messageID)
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
