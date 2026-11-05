package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestNewJSONRespectsLevel(t *testing.T) {
	var buf bytes.Buffer
	l, err := New(&buf, "warn", "json")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	l.Info("ignored")
	l.Warn("kept", "key", "value")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("expected a single JSON entry, got %q: %v", buf.String(), err)
	}
	if entry["msg"] != "kept" || entry["key"] != "value" {
		t.Errorf("unexpected entry: %v", entry)
	}
}

func TestNewRejectsInvalidInput(t *testing.T) {
	if _, err := New(&bytes.Buffer{}, "loud", "json"); err == nil {
		t.Error("expected error for invalid level")
	}
	if _, err := New(&bytes.Buffer{}, "info", "xml"); err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestContextRoundTrip(t *testing.T) {
	if FromContext(context.Background()) != slog.Default() {
		t.Error("FromContext() without logger should return slog.Default()")
	}

	l := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	if FromContext(WithContext(context.Background(), l)) != l {
		t.Error("FromContext() did not return the stored logger")
	}
}
