package repository

import (
	"context"
	"time"
	db "user-management/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, name string, dob time.Time) (db.User, error)
	GetUserByID(ctx context.Context, id int32) (db.User, error)
	UpdateUser(ctx context.Context, id int32, name string, dob time.Time) (db.User, error)
	DeleteUser(ctx context.Context, id int32) error
	ListUsers(ctx context.Context) ([]db.User, error)
	ListUsersWithPagination(ctx context.Context, limit, offset int32) ([]db.User, error)
}

type PostgresUserRepository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) UserRepository {
	return &PostgresUserRepository{
		queries: db.New(pool),
		pool:    pool,
	}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, name string, dob time.Time) (db.User, error) {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		Name: name,
		Dob:  dob,
	})
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id int32) (db.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

func (r *PostgresUserRepository) UpdateUser(ctx context.Context, id int32, name string, dob time.Time) (db.User, error) {
	return r.queries.UpdateUser(ctx, db.UpdateUserParams{
		Name: name,
		Dob:  dob,
		ID:   id,
	})
}

func (r *PostgresUserRepository) DeleteUser(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *PostgresUserRepository) ListUsers(ctx context.Context) ([]db.User, error) {
	return r.queries.ListUsers(ctx)
}

func (r *PostgresUserRepository) ListUsersWithPagination(ctx context.Context, limit, offset int32) ([]db.User, error) {
	return r.queries.ListUsersWithPagination(ctx, db.ListUsersWithPaginationParams{
		Limit:  limit,
		Offset: offset,
	})
}
