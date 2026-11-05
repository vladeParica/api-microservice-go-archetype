package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/vladeParica/api-microservice-go-archetype/internal/httpx"
	"github.com/vladeParica/api-microservice-go-archetype/internal/logger"
)

// Handler exposes the user Service over HTTP.
type Handler struct {
	svc *Service
}

// NewHandler returns a Handler backed by svc.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers the user endpoints on mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/users", h.create)
	mux.HandleFunc("GET /api/v1/users/{id}", h.get)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var p CreateParams
	if err := httpx.DecodeJSON(w, r, &p); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	u, err := h.svc.Create(r.Context(), p)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	w.Header().Set("Location", "/api/v1/users/"+strconv.FormatInt(u.ID, 10))
	httpx.WriteJSON(w, r, http.StatusCreated, u)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		httpx.WriteError(w, r, http.StatusBadRequest, "id must be a positive integer")
		return
	}

	u, err := h.svc.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, u)
}

// writeError maps domain errors to HTTP responses.
func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	var vErr *ValidationError
	switch {
	case errors.As(err, &vErr):
		httpx.WriteError(w, r, http.StatusBadRequest, vErr.Error())
	case errors.Is(err, ErrNotFound):
		httpx.WriteError(w, r, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrEmailTaken):
		httpx.WriteError(w, r, http.StatusConflict, err.Error())
	default:
		logger.FromContext(r.Context()).Error("user request failed", "error", err)
		httpx.WriteError(w, r, http.StatusInternalServerError, "internal server error")
	}
}
