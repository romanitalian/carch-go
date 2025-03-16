package grpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/romanitalian/carch-go/internal/domain"
	"github.com/romanitalian/carch-go/internal/pkg/logger"
	"github.com/romanitalian/carch-go/internal/service"
)

// Mock for UserService
type MockUserService struct {
	mock.Mock
	service.UserService
}

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

// Helper function to create a buffered listener for gRPC testing
func newBufferedListener() *bufconn.Listener {
	return bufconn.Listen(1024 * 1024)
}

// Helper function to create a gRPC connection through a buffered listener
func dialBufferedGrpc(ctx context.Context, listener *bufconn.Listener) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithInsecure(),
	)
}

func TestServer_Run(t *testing.T) {
	// Arrange
	log := logger.New()

	services := &service.Services{
		User: &service.UserService{},
		Log:  log,
	}

	listener := newBufferedListener()
	server := NewServer("bufnet", services, log)

	// Act & Assert
	go func() {
		err := server.Run(listener)
		assert.NoError(t, err)
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Clean up
	server.Shutdown(context.Background())
}

func TestServer_Shutdown(t *testing.T) {
	// Arrange
	log := logger.New()

	services := &service.Services{
		User: &service.UserService{},
		Log:  log,
	}

	listener := newBufferedListener()
	server := NewServer("bufnet", services, log)

	// Act & Assert
	go func() {
		err := server.Run(listener)
		assert.NoError(t, err)
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Shutdown
	err := server.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestServer_Shutdown_WithTimeout(t *testing.T) {
	// Arrange
	log := logger.New()

	services := &service.Services{
		User: &service.UserService{},
		Log:  log,
	}

	listener := newBufferedListener()
	server := NewServer("bufnet", services, log)

	// Act & Assert
	go func() {
		err := server.Run(listener)
		assert.NoError(t, err)
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Shutdown with timeout
	err := server.Shutdown(ctx)
	assert.NoError(t, err)
}
