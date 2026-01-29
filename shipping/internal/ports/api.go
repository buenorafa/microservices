package ports

import "github.com/buenorafa/microservices/shipping/internal/application/core/domain"

type ShippingPort interface {
	CalculateDelivery(shipping domain.Shipping) (int32, error)
}