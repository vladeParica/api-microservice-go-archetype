// Package user implements the user resource: domain model, business rules,
// persistence and HTTP transport.
package user

import (
	"errors"
	"time"
)

// User is the domain representation of a user.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// CreateParams holds the data required to create a user.
type CreateParams struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	// ErrNotFound is returned when the requested user does not exist.
	ErrNotFound = errors.New("user not found")
	// ErrEmailTaken is returned when the email is already registered.
	ErrEmailTaken = errors.New("email already registered")
)

// ValidationError describes invalid input supplied by the caller.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
