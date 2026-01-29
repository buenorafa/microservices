package grpc

import (
	"context"

	shippingpb "github.com/buenorafa/microservices-proto/golang/shipping"
	"github.com/buenorafa/microservices/shipping/internal/application/core/domain"
	"github.com/buenorafa/microservices/shipping/internal/ports"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	shippingpb.UnimplementedShippingServer
	api ports.ShippingPort
}

func NewAdapter(api ports.ShippingPort) *Adapter {
	return &Adapter{api: api}
}

func (a *Adapter) Create(ctx context.Context, req *shippingpb.CreateShippingRequest) (*shippingpb.CreateShippingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	items := make([]domain.ShippingItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, domain.ShippingItem{Quantity: it.Quantity})
	}

	ship := domain.Shipping{
		OrderID: req.OrderId,
		Items:   items,
	}

	days, err := a.api.CalculateDelivery(ship)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &shippingpb.CreateShippingResponse{
		DeliveryDays: int64(days),
	}, nil
} 