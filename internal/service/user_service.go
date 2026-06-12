package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	db "user-management/db/sqlc"
	"user-management/internal/models"
	"user-management/internal/repository"

	"github.com/jackc/pgx/v5"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error)
	GetUserByID(ctx context.Context, id int32) (models.UserWithAgeResponse, error)
	UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error)
	DeleteUser(ctx context.Context, id int32) error
	ListUsers(ctx context.Context, page, limit int) ([]models.UserWithAgeResponse, error)
}

type DefaultUserService struct {
	repo repository.UserRepository
}

func NewDefaultUserService(repo repository.UserRepository) UserService {
	return &DefaultUserService{repo: repo}
}

func (s *DefaultUserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("invalid dob format: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, req.Name, dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format("2006-01-02"),
	}, nil
}

func (s *DefaultUserService) GetUserByID(ctx context.Context, id int32) (models.UserWithAgeResponse, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserWithAgeResponse{}, ErrUserNotFound
		}
		return models.UserWithAgeResponse{}, err
	}

	return models.UserWithAgeResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format("2006-01-02"),
		Age:  CalculateAge(user.Dob),
	}, nil
}

func (s *DefaultUserService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error) {
	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("invalid dob format: %w", err)
	}

	user, err := s.repo.UpdateUser(ctx, id, req.Name, dob)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserResponse{}, ErrUserNotFound
		}
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Format("2006-01-02"),
	}, nil
}

func (s *DefaultUserService) DeleteUser(ctx context.Context, id int32) error {
	_, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	return s.repo.DeleteUser(ctx, id)
}

func (s *DefaultUserService) ListUsers(ctx context.Context, page, limit int) ([]models.UserWithAgeResponse, error) {
	var users []db.User
	var err error

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		users, err = s.repo.ListUsersWithPagination(ctx, int32(limit), int32(offset))
	} else {
		users, err = s.repo.ListUsers(ctx)
	}

	if err != nil {
		return nil, err
	}

	res := make([]models.UserWithAgeResponse, len(users))
	for i, u := range users {
		res[i] = models.UserWithAgeResponse{
			ID:   u.ID,
			Name: u.Name,
			DOB:  u.Dob.Format("2006-01-02"),
			Age:  CalculateAge(u.Dob),
		}
	}
	return res, nil
}

func CalculateAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		years--
	}
	return years
}
