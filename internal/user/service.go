package user

import (
	"context"
	"net/mail"
	"strings"
)

// Service implements the user use cases.
type Service struct {
	repo Repository
}

// NewService returns a Service that stores users in repo.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create validates p and stores a new user.
func (s *Service) Create(ctx context.Context, p CreateParams) (User, error) {
	p.Name = strings.TrimSpace(p.Name)
	p.Email = strings.ToLower(strings.TrimSpace(p.Email))

	if p.Name == "" {
		return User{}, &ValidationError{Field: "name", Message: "is required"}
	}
	if addr, err := mail.ParseAddress(p.Email); err != nil || addr.Address != p.Email {
		return User{}, &ValidationError{Field: "email", Message: "must be a valid email address"}
	}

	return s.repo.Create(ctx, p)
}

// Get returns the user with the given id.
func (s *Service) Get(ctx context.Context, id int64) (User, error) {
	return s.repo.GetByID(ctx, id)
}
