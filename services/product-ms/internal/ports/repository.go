package ports

import (
	"context"
	"product-microservice/internal/domain"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, p *domain.Product) (*domain.Product, error)
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	FindAll(ctx context.Context) ([]domain.Product, error)
	Update(ctx context.Context, p *domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id string) error
}