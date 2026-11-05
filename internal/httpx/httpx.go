// Package httpx contains small helpers shared by HTTP handlers.
package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/vladeParica/api-microservice-go-archetype/internal/logger"
)

// maxBodyBytes limits the size of JSON request bodies.
const maxBodyBytes = 1 << 20 // 1 MiB

// ErrorResponse is the body returned for every non-2xx response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteJSON encodes v as JSON and writes it with the given status code.
func WriteJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.FromContext(r.Context()).Error("encode response", "error", err)
	}
}

// WriteError writes an ErrorResponse with the given status code and message.
func WriteError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	WriteJSON(w, r, status, ErrorResponse{Error: msg})
}

// DecodeJSON decodes the request body into dst, rejecting unknown fields,
// trailing data and bodies larger than 1 MiB.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("invalid JSON body: must contain a single object")
	}
	return nil
}
