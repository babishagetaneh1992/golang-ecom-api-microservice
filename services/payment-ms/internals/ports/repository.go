package ports

import (
	"context"
	"payment-microservice/internals/domain"
)

type PaymentRepository interface {
	Create(ctx context.Context, p *domain.Payment) (*domain.Payment, error)
	FindByID(ctx context.Context, id string) (*domain.Payment, error)
	List(ctx context.Context) ([]*domain.Payment, error)
	UpdateStatus(ctx context.Context, id string, status string) (*domain.Payment, error)
	Delete(ctx context.Context, id string) error
}