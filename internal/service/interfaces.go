package service

import (
	"context"

	"github.com/romanitalian/carch-go/internal/domain"
)

// UserServiceInterface defines the interface for user service
type UserServiceInterface interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.User, error)
}
