package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/vladeParica/api-microservice-go-archetype/internal/database"
)

// pgUniqueViolation is the PostgreSQL error code for unique_violation.
const pgUniqueViolation = "23505"

// Repository persists users.
type Repository interface {
	Create(ctx context.Context, p CreateParams) (User, error)
	GetByID(ctx context.Context, id int64) (User, error)
}

// PostgresRepository is a Repository backed by sqlc generated queries.
type PostgresRepository struct {
	q database.Querier
}

// NewPostgresRepository returns a Repository that uses q.
func NewPostgresRepository(q database.Querier) *PostgresRepository {
	return &PostgresRepository{q: q}
}

// Create inserts a new user.
func (r *PostgresRepository) Create(ctx context.Context, p CreateParams) (User, error) {
	row, err := r.q.CreateUser(ctx, database.CreateUserParams{Name: p.Name, Email: p.Email})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return User{}, ErrEmailTaken
		}
		return User{}, fmt.Errorf("insert user: %w", err)
	}
	return fromRow(row), nil
}

// GetByID returns the user with the given id.
func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("select user %d: %w", id, err)
	}
	return fromRow(row), nil
}

func fromRow(row database.User) User {
	return User{ID: row.ID, Name: row.Name, Email: row.Email, CreatedAt: row.CreatedAt}
}
