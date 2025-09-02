package ports

import "user-microservice/internal/domain"

// Outbound port (persistent storage)
type UserRepository interface {
	Create(user *domain.User) error
	GetById(id string) (*domain.User, error)
	GetAll() ([]domain.User, error)
	Exists(id string) (bool, error)
	FindByEmail(email string) (*domain.User, error)
}

