package repository

import (
	"github.com/romanitalian/carch-go/internal/domain"
)

type Repositories struct {
	User domain.UserRepository
}

// NewRepositories creates a new Repositories instance
func NewRepositories(db *DB, mq *RabbitMQ) *Repositories {
	return &Repositories{
		User: NewUserRepository(db.DB),
	}
}
