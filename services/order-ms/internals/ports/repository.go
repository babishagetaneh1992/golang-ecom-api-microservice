package ports

import (
	"context"
	"order-microservice/internals/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order)(*domain.Order, error)
	FindByID(ctx context.Context, id string)(*domain.Order, error)
	List(ctx context.Context) ([]*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status string) (*domain.Order, error)
	Delete(ctx context.Context, id string) error
}