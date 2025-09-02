package ports

import (
	"context"
	"product-microservice/internal/domain"
)

type ProductService interface {
	CreateNewProduct(ctx context.Context, p *domain.Product) (*domain.Product, error) 
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]domain.Product, error)
	UpdateProduct(ctx context.Context, p *domain.Product) (*domain.Product,error)
	DeleteProduct(ctx context.Context, id string) error
}