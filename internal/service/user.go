package service

import (
	"context"

	"github.com/romanitalian/carch-go/internal/domain"
	"github.com/romanitalian/carch-go/internal/pkg/logger"
)

type UserService struct {
	repo domain.UserRepository
	log  *logger.Logger
}

func NewUserService(repo domain.UserRepository, log *logger.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
	}
}

func (s *UserService) Create(ctx context.Context, user *domain.User) error {
	// Business logic and validation
	s.log.Info("Creating user", map[string]interface{}{"user_id": user.ID})
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	s.log.Info("Getting user by ID", map[string]interface{}{"user_id": id})
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, user *domain.User) error {
	s.log.Info("Updating user", map[string]interface{}{"user_id": user.ID})
	return s.repo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	s.log.Info("Deleting user", map[string]interface{}{"user_id": id})
	return s.repo.Delete(ctx, id)
}

func (s *UserService) List(ctx context.Context) ([]*domain.User, error) {
	s.log.Info("Listing users", nil)
	return s.repo.List(ctx)
}
