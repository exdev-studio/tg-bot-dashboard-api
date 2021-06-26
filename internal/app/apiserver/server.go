package apiserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ctxKey int8

type server struct {
	router *mux.Router
	logger *logrus.Logger
}

const (
	ctxKeyRequestID ctxKey = iota
)

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer(logger *logrus.Logger) *server {
	s := &server{
		logger: logger,
		router: mux.NewRouter(),
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(s.setContentType)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	s.router.HandleFunc("/", s.handleDefaultRequest()).Methods(http.MethodPost)
}

func (s *server) setContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *server) handleDefaultRequest() http.HandlerFunc {
	type chat struct {
		Id       int64  `json:"id"`
		Type     string `json:"type"`
		Username string `json:"username"`
	}
	type message struct {
		Id   int    `json:"message_id"`
		Text string `json:"text"`
		Chat chat   `json:"chat"`
	}

	type request struct {
		UpdateId int     `json:"update_id"`
		Message  message `json:"message"`
		Chat     chat    `json:"chat"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}

		s.logger.WithFields(logrus.Fields{
			"update-id":     req.UpdateId,
			"message-id":    req.Message.Id,
			"message":       req.Message.Text,
			"chat-id":       req.Message.Chat.Id,
			"chat-type":     req.Message.Chat.Type,
			"chat-username": req.Message.Chat.Username,
		}).Debug("request received")

		s.respond(w, http.StatusOK, http.StatusText(http.StatusOK))
	}
}

func (s *server) error(w http.ResponseWriter, code int, err error) {
	s.respond(w, code, map[string]string{
		"error": err.Error(),
	})
}

func (s *server) respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			s.logger.Fatal(err)
		}
	}
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote-addr": r.RemoteAddr,
			"request-id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Debugf("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Debugf(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}
