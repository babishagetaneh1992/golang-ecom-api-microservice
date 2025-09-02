package ports

import "user-microservice/internal/domain"


// Inbound port (use cases)
type UserService interface {
	Register(user *domain.User) (*domain.User, error)
	GetUser(id string) (*domain.User, error)
	ListUsers() ([]domain.User, error)
	Exists(id string) (bool, error)
	Authenticate(email, password string)(*domain.User, error)
}

