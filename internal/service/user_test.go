package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/romanitalian/carch-go/internal/domain"
	"github.com/romanitalian/carch-go/internal/pkg/logger"
)

// Мок для UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func TestUserService_Create(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	log := logger.New()
	service := NewUserService(mockRepo, log)
	ctx := context.Background()

	user := &domain.User{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	// Настройка мока
	mockRepo.On("Create", ctx, user).Return(nil)

	// Act
	err := service.Create(ctx, user)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	log := logger.New()
	service := NewUserService(mockRepo, log)
	ctx := context.Background()

	userID := "user-123"
	expectedUser := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Настройка мока
	mockRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)

	// Act
	user, err := service.GetByID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	log := logger.New()
	service := NewUserService(mockRepo, log)
	ctx := context.Background()

	userID := "non-existent-user"

	// Настройка мока
	mockRepo.On("GetByID", ctx, userID).Return(nil, domain.ErrUserNotFound)

	// Act
	user, err := service.GetByID(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Update(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	log := logger.New()
	service := NewUserService(mockRepo, log)
	ctx := context.Background()

	user := &domain.User{
		ID:    "user-123",
		Email: "updated@example.com",
		Name:  "Updated User",
	}

	// Настройка мока
	mockRepo.On("Update", ctx, user).Return(nil)

	// Act
	err := service.Update(ctx, user)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	log := logger.New()
	service := NewUserService(mockRepo, log)
	ctx := context.Background()

	user := &domain.User{
		ID:    "non-existent-user",
		Email: "updated@example.com",
		Name:  "Updated User",
	}

	// Настройка мока
	mockRepo.On("Update", ctx, user).Return(domain.ErrUserNotFound)

	// Act
	err := service.Update(ctx, user)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	log := logger.New()
	service := NewUserService(mockRepo, log)
	ctx := context.Background()

	userID := "user-123"

	// Настройка мока
	mockRepo.On("Delete", ctx, userID).Return(nil)

	// Act
	err := service.Delete(ctx, userID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	log := logger.New()
	service := NewUserService(mockRepo, log)
	ctx := context.Background()

	userID := "non-existent-user"

	// Настройка мока
	mockRepo.On("Delete", ctx, userID).Return(domain.ErrUserNotFound)

	// Act
	err := service.Delete(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_List(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	log := logger.New()
	service := NewUserService(mockRepo, log)
	ctx := context.Background()

	expectedUsers := []*domain.User{
		{
			ID:        "user-1",
			Email:     "user1@example.com",
			Name:      "User 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "user-2",
			Email:     "user2@example.com",
			Name:      "User 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Настройка мока
	mockRepo.On("List", ctx).Return(expectedUsers, nil)

	// Act
	users, err := service.List(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	assert.Len(t, users, 2)
	mockRepo.AssertExpectations(t)
}

func TestUserService_List_Error(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	log := logger.New()
	service := NewUserService(mockRepo, log)
	ctx := context.Background()

	expectedError := errors.New("database error")

	// Настройка мока
	mockRepo.On("List", ctx).Return(nil, expectedError)

	// Act
	users, err := service.List(ctx)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, users)
	mockRepo.AssertExpectations(t)
}
