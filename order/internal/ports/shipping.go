package ports

import "github.com/buenorafa/microservices/order/internal/application/core/domain"

type ShippingPort interface {
	RequestShipping(order *domain.Order) (int32, error) // retorna deliveryDays
}