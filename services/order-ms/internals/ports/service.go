package ports

import (
	"context"
	"order-microservice/internals/domain"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	GetOrder(ctx context.Context, id string) (*domain.Order, error)
	ListOrders(ctx context.Context) ([]*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status string) (*domain.Order, error)
	DeleteOrder(ctx context.Context, id string) error
	CreateOrderFromCart(ctx context.Context, userID string) (*domain.Order, error)
}
