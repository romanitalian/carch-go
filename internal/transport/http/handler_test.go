package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/romanitalian/carch-go/internal/domain"
	"github.com/romanitalian/carch-go/internal/pkg/logger"
	"github.com/romanitalian/carch-go/internal/service"
)

// Mock for UserService
type MockUserService struct {
	mock.Mock
}

// Ensure MockUserService implements service.UserServiceInterface
var _ service.UserServiceInterface = (*MockUserService)(nil)

func (m *MockUserService) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) List(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

// Helper function to set up test environment
func setupTestHandler() (*MockUserService, *Handler, *http.ServeMux) {
	mockUserService := new(MockUserService)
	log := logger.New()

	// Create a service.Services struct with our mock
	services := &service.Services{
		User: mockUserService,
		Log:  log,
	}

	handler := NewHandler(services, log)

	return mockUserService, handler, handler.mux
}

func TestHandler_createUser(t *testing.T) {
	// Arrange
	mockUserService, handler, _ := setupTestHandler()

	// Prepare request
	reqBody := createUserRQ{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Mock service behavior
	mockUserService.On("Create", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.Email == reqBody.Email && user.Password == reqBody.Password && user.Name == reqBody.Name
	})).Return(nil)

	// Act
	handler.createUser(rr, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rr.Code)
	mockUserService.AssertExpectations(t)
}

func TestHandler_createUser_ValidationError(t *testing.T) {
	// Arrange
	mockUserService, handler, _ := setupTestHandler()

	// Prepare invalid request (missing required fields)
	reqBody := `{"email": "invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// We don't expect the service to be called, but we need to mock it anyway
	// because the handler will try to call it if the JSON parsing succeeds
	mockUserService.On("Create", mock.Anything, mock.Anything).Return(nil).Maybe()

	// Act
	handler.createUser(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_createUser_ServiceError(t *testing.T) {
	// Arrange
	mockUserService, handler, _ := setupTestHandler()

	// Prepare request
	reqBody := createUserRQ{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Mock service error
	expectedErr := errors.New("service error")
	mockUserService.On("Create", mock.Anything, mock.Anything).Return(expectedErr)

	// Act
	handler.createUser(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockUserService.AssertExpectations(t)
}

func TestHandler_getUserByID(t *testing.T) {
	// Arrange
	mockUserService, handler, _ := setupTestHandler()

	userID := "user-123"
	expectedUser := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create request with path parameter
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+userID, nil)

	// Mock PathValue to return the ID
	origPathValueFunc := pathValueFunc
	defer func() { pathValueFunc = origPathValueFunc }()
	pathValueFunc = func(r *http.Request, key string) string {
		if key == "id" {
			return userID
		}
		return ""
	}

	rr := httptest.NewRecorder()

	// Mock service behavior
	mockUserService.On("GetByID", mock.Anything, userID).Return(expectedUser, nil)

	// Act
	handler.getUserByID(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)

	var responseUser domain.User
	err := json.Unmarshal(rr.Body.Bytes(), &responseUser)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, responseUser.ID)
	assert.Equal(t, expectedUser.Email, responseUser.Email)
	assert.Equal(t, expectedUser.Name, responseUser.Name)

	mockUserService.AssertExpectations(t)
}

func TestHandler_getUserByID_NotFound(t *testing.T) {
	// Arrange
	mockUserService, handler, _ := setupTestHandler()

	userID := "non-existent-id"

	// Create request with path parameter
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+userID, nil)

	// Mock PathValue to return the ID
	origPathValueFunc := pathValueFunc
	defer func() { pathValueFunc = origPathValueFunc }()
	pathValueFunc = func(r *http.Request, key string) string {
		if key == "id" {
			return userID
		}
		return ""
	}

	rr := httptest.NewRecorder()

	// Mock service behavior
	mockUserService.On("GetByID", mock.Anything, userID).Return(nil, domain.ErrUserNotFound)

	// Act
	handler.getUserByID(rr, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockUserService.AssertExpectations(t)
}

func TestHandler_updateUser(t *testing.T) {
	// Arrange
	mockUserService, handler, _ := setupTestHandler()

	userID := "user-123"
	reqBody := updateUserRQ{
		Email: "updated@example.com",
		Name:  "Updated User",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/"+userID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Mock PathValue to return the ID
	origPathValueFunc := pathValueFunc
	defer func() { pathValueFunc = origPathValueFunc }()
	pathValueFunc = func(r *http.Request, key string) string {
		if key == "id" {
			return userID
		}
		return ""
	}

	rr := httptest.NewRecorder()

	// Mock service behavior
	mockUserService.On("Update", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.ID == userID && user.Email == reqBody.Email && user.Name == reqBody.Name
	})).Return(nil)

	// Act
	handler.updateUser(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	mockUserService.AssertExpectations(t)
}

func TestHandler_deleteUser(t *testing.T) {
	// Arrange
	mockUserService, handler, _ := setupTestHandler()

	userID := "user-123"

	// Create request with path parameter
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+userID, nil)

	// Mock PathValue to return the ID
	origPathValueFunc := pathValueFunc
	defer func() { pathValueFunc = origPathValueFunc }()
	pathValueFunc = func(r *http.Request, key string) string {
		if key == "id" {
			return userID
		}
		return ""
	}

	rr := httptest.NewRecorder()

	// Mock service behavior
	mockUserService.On("Delete", mock.Anything, userID).Return(nil)

	// Act
	handler.deleteUser(rr, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockUserService.AssertExpectations(t)
}

func TestHandler_listUsers(t *testing.T) {
	// Arrange
	mockUserService, handler, _ := setupTestHandler()

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

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rr := httptest.NewRecorder()

	// Mock service behavior
	mockUserService.On("List", mock.Anything).Return(expectedUsers, nil)

	// Act
	handler.listUsers(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)

	var responseUsers []*domain.User
	err := json.Unmarshal(rr.Body.Bytes(), &responseUsers)
	assert.NoError(t, err)
	assert.Len(t, responseUsers, len(expectedUsers))

	mockUserService.AssertExpectations(t)
}
