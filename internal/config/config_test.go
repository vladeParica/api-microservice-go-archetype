package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadDefaultsWhenFileIsMissing(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.HTTP.Port != 8080 {
		t.Errorf("HTTP.Port = %d, want 8080", cfg.HTTP.Port)
	}
	if cfg.HTTP.ShutdownTimeout != 10*time.Second {
		t.Errorf("HTTP.ShutdownTimeout = %s, want 10s", cfg.HTTP.ShutdownTimeout)
	}
}

func TestLoadFileAndEnvOverride(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	content := []byte("http:\n  port: 9000\nlog:\n  level: debug\ndatabase:\n  host: db\n")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DATABASE_HOST", "from-env")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.HTTP.Port != 9000 {
		t.Errorf("HTTP.Port = %d, want 9000", cfg.HTTP.Port)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("Log.Level = %q, want debug", cfg.Log.Level)
	}
	if cfg.Database.Host != "from-env" {
		t.Errorf("Database.Host = %q, want from-env", cfg.Database.Host)
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	t.Setenv("LOG_LEVEL", "verbose")

	if _, err := Load(""); err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}
}

func TestDatabaseURLEscapesCredentials(t *testing.T) {
	db := Database{Host: "localhost", Port: 5432, User: "app", Password: "p@ss/word", Name: "app", SSLMode: "disable"}

	want := "postgres://app:p%40ss%2Fword@localhost:5432/app?sslmode=disable"
	if got := db.URL(); got != want {
		t.Errorf("URL() = %q, want %q", got, want)
	}
}
