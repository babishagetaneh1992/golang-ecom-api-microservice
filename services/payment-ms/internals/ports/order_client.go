package ports

import (
	"context"
	"order-microservice/adaptors/grpc/pb/order-microservice/services/order-ms/adaptors/grpc/pb"
	//"payment-microservice/internals/domain"
)

type OrderClient interface {
    UpdateOrderStatus(ctx context.Context, orderID, status string) error
	GetOrder(ctx context.Context, orderID string) (*pb.Order, error)
}