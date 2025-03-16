package service

import (
	"github.com/romanitalian/carch-go/internal/domain"
	"github.com/romanitalian/carch-go/internal/pkg/logger"
)

type Deps struct {
	Repos        *Repositories
	MessageQueue interface{}
	Logger       *logger.Logger
}

type Repositories struct {
	User domain.UserRepository
}

type Services struct {
	User UserServiceInterface
	Log  *logger.Logger
}

func NewServices(deps Deps) *Services {
	return &Services{
		User: NewUserService(deps.Repos.User, deps.Logger),
		Log:  deps.Logger,
	}
}
