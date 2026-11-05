package httpx

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSON(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{name: "valid", body: `{"name":"john"}`},
		{name: "unknown field", body: `{"name":"john","age":3}`, wantErr: true},
		{name: "trailing data", body: `{"name":"john"}{}`, wantErr: true},
		{name: "malformed", body: `{"name":`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(tt.body))
			var dst struct {
				Name string `json:"name"`
			}

			err := DecodeJSON(httptest.NewRecorder(), r, &dst)
			if (err != nil) != tt.wantErr {
				t.Fatalf("DecodeJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	WriteError(w, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil), http.StatusNotFound, "not found")

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	if got, want := strings.TrimSpace(w.Body.String()), `{"error":"not found"}`; got != want {
		t.Errorf("body = %s, want %s", got, want)
	}
}
