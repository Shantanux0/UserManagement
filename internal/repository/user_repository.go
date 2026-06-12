package repository

import (
	"context"
	"time"
	db "user-management/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository defines the interface for database operations on users.
type UserRepository interface {
	CreateUser(ctx context.Context, name string, dob time.Time) (db.User, error)
	GetUserByID(ctx context.Context, id int32) (db.User, error)
	UpdateUser(ctx context.Context, id int32, name string, dob time.Time) (db.User, error)
	DeleteUser(ctx context.Context, id int32) error
	ListUsers(ctx context.Context) ([]db.User, error)
	ListUsersWithPagination(ctx context.Context, limit, offset int32) ([]db.User, error)
}

// PostgresUserRepository implements UserRepository using pgx and SQLC.
type PostgresUserRepository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

// NewPostgresUserRepository creates a new PostgresUserRepository.
func NewPostgresUserRepository(pool *pgxpool.Pool) UserRepository {
	return &PostgresUserRepository{
		queries: db.New(pool),
		pool:    pool,
	}
}

// CreateUser inserts a user into the database.
func (r *PostgresUserRepository) CreateUser(ctx context.Context, name string, dob time.Time) (db.User, error) {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		Name: name,
		Dob:  dob,
	})
}

// GetUserByID fetches a user by their ID.
func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id int32) (db.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

// UpdateUser updates an existing user's details.
func (r *PostgresUserRepository) UpdateUser(ctx context.Context, id int32, name string, dob time.Time) (db.User, error) {
	return r.queries.UpdateUser(ctx, db.UpdateUserParams{
		Name: name,
		Dob:  dob,
		ID:   id,
	})
}

// DeleteUser deletes a user by ID.
func (r *PostgresUserRepository) DeleteUser(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}

// ListUsers lists all users in the database.
func (r *PostgresUserRepository) ListUsers(ctx context.Context) ([]db.User, error) {
	return r.queries.ListUsers(ctx)
}

// ListUsersWithPagination lists users with pagination support.
func (r *PostgresUserRepository) ListUsersWithPagination(ctx context.Context, limit, offset int32) ([]db.User, error) {
	return r.queries.ListUsersWithPagination(ctx, db.ListUsersWithPaginationParams{
		Limit:  limit,
		Offset: offset,
	})
}
