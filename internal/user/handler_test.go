package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// fakeRepository is an in-memory Repository for tests.
type fakeRepository struct {
	mu     sync.Mutex
	users  map[int64]User
	nextID int64
	err    error
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{users: map[int64]User{}}
}

func (f *fakeRepository) Create(_ context.Context, p CreateParams) (User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		return User{}, f.err
	}
	for _, u := range f.users {
		if u.Email == p.Email {
			return User{}, ErrEmailTaken
		}
	}
	f.nextID++
	u := User{ID: f.nextID, Name: p.Name, Email: p.Email, CreatedAt: time.Now().UTC()}
	f.users[u.ID] = u
	return u, nil
}

func (f *fakeRepository) GetByID(_ context.Context, id int64) (User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		return User{}, f.err
	}
	u, ok := f.users[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

func newTestMux(repo Repository) *http.ServeMux {
	mux := http.NewServeMux()
	NewHandler(NewService(repo)).RegisterRoutes(mux)
	return mux
}

func do(t *testing.T, h http.Handler, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequestWithContext(t.Context(), method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		repoErr    error
		wantStatus int
	}{
		{name: "created", body: `{"name":"John","email":"John@Example.com"}`, wantStatus: http.StatusCreated},
		{name: "missing name", body: `{"name":" ","email":"john@example.com"}`, wantStatus: http.StatusBadRequest},
		{name: "invalid email", body: `{"name":"John","email":"nope"}`, wantStatus: http.StatusBadRequest},
		{name: "unknown field", body: `{"name":"John","email":"john@example.com","password":"x"}`, wantStatus: http.StatusBadRequest},
		{name: "repository failure", body: `{"name":"John","email":"john@example.com"}`, repoErr: errors.New("boom"), wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepository()
			repo.err = tt.repoErr

			rec := do(t, newTestMux(repo), http.MethodPost, "/api/v1/users", tt.body)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body)
			}

			if tt.wantStatus == http.StatusCreated {
				var got User
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}
				if got.ID != 1 || got.Email != "john@example.com" {
					t.Errorf("unexpected user: %+v", got)
				}
				if loc := rec.Header().Get("Location"); loc != "/api/v1/users/1" {
					t.Errorf("Location = %q, want /api/v1/users/1", loc)
				}
			}
		})
	}
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	mux := newTestMux(newFakeRepository())
	body := `{"name":"John","email":"john@example.com"}`

	do(t, mux, http.MethodPost, "/api/v1/users", body)
	rec := do(t, mux, http.MethodPost, "/api/v1/users", body)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestGetUser(t *testing.T) {
	mux := newTestMux(newFakeRepository())
	do(t, mux, http.MethodPost, "/api/v1/users", `{"name":"John","email":"john@example.com"}`)

	tests := []struct {
		name       string
		target     string
		wantStatus int
	}{
		{name: "found", target: "/api/v1/users/1", wantStatus: http.StatusOK},
		{name: "not found", target: "/api/v1/users/42", wantStatus: http.StatusNotFound},
		{name: "invalid id", target: "/api/v1/users/abc", wantStatus: http.StatusBadRequest},
		{name: "non positive id", target: "/api/v1/users/0", wantStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := do(t, mux, http.MethodGet, tt.target, "")
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body)
			}
		})
	}
}
