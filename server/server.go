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

func (s *Server) publish() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		var publishReq kitsune.PublishRequest
		if err := receiveJSON(res, req, &publishReq); err != nil {
			zap.S().Errorw("Error marshalling request", "error", err)
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
			case repository.ErrDuplicateMessage:
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

		message, err := s.repo.GetMessage(params.ByName("topic"), params.ByName("messageId"))
		if err != nil {
			zap.S().Errorw("Error retrieving message", "error", err)
			switch err {
			case repository.ErrMessageNotFound:
				http.Error(res, "Message not found", http.StatusNotFound)
			default:
				http.Error(res, "Not Found", http.StatusNotFound)
			}
			return
		}

		if err := sendJSON(res, &message); err != nil {
			zap.S().Errorw("Error marshalling response", "error", err)
			return
		}
	}
}

func (s *Server) poll() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		topic := httprouter.ParamsFromContext(req.Context()).ByName("topic")

		var pollReq kitsune.PollRequest
		if err := receiveJSON(res, req, &pollReq); err != nil {
			zap.S().Errorw("Error marshalling request", "error", err)
			return
		}

		messages, err := s.repo.PollTopic(topic, pollReq)
		if err != nil {
			zap.S().Errorw("Error polling messages", "error", err)
			switch err {
			case repository.ErrTopicNotFound:
				http.Error(res, "Topic not found", http.StatusNotFound)
			default:
				http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			}

			return
		}

		if err := sendJSON(res, &messages); err != nil {
			zap.S().Errorw("Error marshalling response", "error", err)
			return
		}

		zap.S().Infow("Messages sucessfully polled",
			"topic", topic,
			"subscription", pollReq.SubscriptionName,
			"messagesReturned", len(messages),
		)
	}
}

func (s *Server) setSubscriptionPosition() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		topic := httprouter.ParamsFromContext(req.Context()).ByName("topic")

		var setSubPosReq kitsune.SubscriptionPositionRequest
		if err := receiveJSON(res, req, &setSubPosReq); err != nil {
			zap.S().Errorw("Error marshalling request", "error", err)
			return
		}

		err := s.repo.SetSubscriptionPosition(topic, setSubPosReq)
		if err != nil {
			zap.S().Errorw("Error setting subscription position", "error", err)
			switch err {
			case repository.ErrTopicNotFound:
				http.Error(res, "Topic not found", http.StatusNotFound)
			default:
				http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			}

			return
		}

		zap.S().Infow("Subscription position set",
			"topic", topic,
			"topic", topic,
			"subscription", setSubPosReq.SubscriptionName,
		)
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
