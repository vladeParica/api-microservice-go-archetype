// Package health exposes liveness and readiness probes.
package health

import (
	"context"
	"net/http"
	"time"

	"github.com/vladeParica/api-microservice-go-archetype/internal/httpx"
	"github.com/vladeParica/api-microservice-go-archetype/internal/logger"
)

const pingTimeout = 2 * time.Second

// Pinger is implemented by dependencies that can report their availability,
// such as *pgxpool.Pool.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Handler serves the health endpoints.
type Handler struct {
	version string
	deps    map[string]Pinger
}

// NewHandler returns a Handler reporting version and checking deps on readiness.
func NewHandler(version string, deps map[string]Pinger) *Handler {
	return &Handler{version: version, deps: deps}
}

// RegisterRoutes registers the health endpoints on mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", h.live)
	mux.HandleFunc("GET /readyz", h.ready)
}

type response struct {
	Status  string            `json:"status"`
	Version string            `json:"version,omitempty"`
	Checks  map[string]string `json:"checks,omitempty"`
}

func (h *Handler) live(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, r, http.StatusOK, response{Status: "ok", Version: h.version})
}

func (h *Handler) ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), pingTimeout)
	defer cancel()

	res := response{Status: "ok", Version: h.version, Checks: make(map[string]string, len(h.deps))}
	status := http.StatusOK

	for name, dep := range h.deps {
		if err := dep.Ping(ctx); err != nil {
			logger.FromContext(ctx).Warn("readiness check failed", "dependency", name, "error", err)
			res.Checks[name] = "unavailable"
			res.Status = "unavailable"
			status = http.StatusServiceUnavailable
			continue
		}
		res.Checks[name] = "ok"
	}

	httpx.WriteJSON(w, r, status, res)
}
