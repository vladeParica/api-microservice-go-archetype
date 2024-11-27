package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"example.com/test/pkg/components/databases/repository"
	"example.com/test/pkg/config"
	"example.com/test/pkg/config/logger"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	router  *mux.Router
	Handler http.Handler
	queries *repository.Queries
	db      *pgxpool.Pool
	log     *zap.Logger
}

func NewServer(
	queries *repository.Queries,
	db *pgxpool.Pool,
	log *zap.Logger,
) *Server {
	return &Server{
		router:  mux.NewRouter(),
		queries: queries,
		db:      db,
		log:     log,
	}
}

func (s *Server) registerHandlers() {
	s.router.HandleFunc("/api/v1/hello", s.helloHandler).Methods("GET")
	s.router.HandleFunc("/api/v1/user/create", s.createUserHandler).Methods("POST")
	s.router.HandleFunc("/api/v1/user/{user-id}", s.getUserHandler).Methods("GET")
}

func (s *Server) ListenAndServe(stopCh <-chan struct{}) {
	s.registerHandlers()
	s.Handler = s.router
	s.log.Info("Starting server...")

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.Configuration.MicroserviceServer, config.Configuration.MicroservicePort), RequestIDMiddleware(s.Handler)); err != nil {
			logger.FromCtx().Fatal("Error starting server", zap.Error(err))
		}
	}()
	s.log.Info("Server started")

	// wait for SIGTERM or SIGINT
	<-stopCh
	s.log.Info("Stopping server...")
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Reset()

		requestID := uuid.New()

		ctx := r.Context()
		const keyTraceIdContext = "requestID"
		ctx = context.WithValue(ctx, keyTraceIdContext, requestID)

		r = r.WithContext(ctx)

		logger.LoggerInstance = logger.LoggerInstance.With(zap.String(keyTraceIdContext, requestID.String()))
		r = r.WithContext(logger.WithCtx(ctx, logger.LoggerInstance))

		next.ServeHTTP(w, r)
	})
}

func (s *Server) JSONResponse(w http.ResponseWriter, responseCode int, result interface{}) {
	body, err := json.Marshal(result)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.log.Error("Error marshalling response", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	if rawData, err := prettyJSON(body); err == nil {
		w.Write(rawData)
	}

}

func prettyJSON(data []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", "  ")
	return out.Bytes(), err
}
