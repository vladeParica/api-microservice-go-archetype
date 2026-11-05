package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type pingerFunc func(ctx context.Context) error

func (f pingerFunc) Ping(ctx context.Context) error { return f(ctx) }

func TestHealthEndpoints(t *testing.T) {
	healthy := pingerFunc(func(context.Context) error { return nil })
	broken := pingerFunc(func(context.Context) error { return errors.New("down") })

	tests := []struct {
		name       string
		target     string
		deps       map[string]Pinger
		wantStatus int
	}{
		{name: "live", target: "/healthz", deps: map[string]Pinger{"db": broken}, wantStatus: http.StatusOK},
		{name: "ready", target: "/readyz", deps: map[string]Pinger{"db": healthy}, wantStatus: http.StatusOK},
		{name: "not ready", target: "/readyz", deps: map[string]Pinger{"db": broken}, wantStatus: http.StatusServiceUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			NewHandler("test", tt.deps).RegisterRoutes(mux)

			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.target, nil))

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body)
			}
		})
	}
}
