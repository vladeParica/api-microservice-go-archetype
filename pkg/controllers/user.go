package controllers

import (
	"encoding/json"
	"github.com/vladeParica/api-microservice-go-archetype/internal/database"
	"github.com/vladeParica/api-microservice-go-archetype/pkg/config/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (s *Server) createUserHandler(response http.ResponseWriter, request *http.Request) {

	var user database.User
	decoder := json.NewDecoder(request.Body)

	if err := decoder.Decode(&user); err != nil {
		logger.FromCtx().Error("Error decoding request body", zap.Error(err))
	}

	s.log.Info("request parsed", zap.Any("user", user))

	createParams := database.CreateUserParams{
		Name:  user.Name,
		Email: user.Email,
	}

	tx, err := s.db.Begin(request.Context())

	if err != nil {
		panic(err)
	}

	qtx := s.queries.WithTx(tx)
	result, err := qtx.CreateUser(request.Context(), createParams)
	if err != nil {
		s.log.Error("Error creating user", zap.Error(err))
		s.JSONResponse(response, http.StatusInternalServerError, map[string]string{"message": "Invalid request"})
		return
	}

	tx.Commit(request.Context())
	s.JSONResponse(response, http.StatusCreated, map[string]string{"message": "User created", "id": strconv.FormatInt(result.ID, 10)})

}

func (s *Server) getUserHandler(response http.ResponseWriter, request *http.Request) {

	idParam := mux.Vars(request)["user-id"]

	id, err := strconv.Atoi(idParam)
	if err != nil {
		s.JSONResponse(response, http.StatusBadRequest, map[string]string{"message": "Invalid 'id' format"})
		return
	}

	tx, err := s.db.Begin(request.Context())
	if err != nil {
		panic(err)
	}

	qtx := s.queries.WithTx(tx)
	result, err := qtx.GetUserByID(request.Context(), int64(id))
	if err != nil {
		s.log.Error("Error consulting user", zap.Error(err))
		s.JSONResponse(response, http.StatusInternalServerError, map[string]string{"message": "Invalid request"})
		return
	}

	tx.Commit(request.Context())
	s.JSONResponse(response, http.StatusOK, result)
}
