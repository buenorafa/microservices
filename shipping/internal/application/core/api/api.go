package api

import (
	"errors"

	"github.com/buenorafa/microservices/shipping/internal/application/core/domain"
	"github.com/buenorafa/microservices/shipping/internal/ports"
)

type Application struct{}

// garante que implementa o port
var _ ports.ShippingPort = (*Application)(nil)

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) CalculateDelivery(shipping domain.Shipping) (int32, error) {

	totalUnits := shipping.TotalUnits()
	if totalUnits <= 0 {
		return 0, errors.New("invalid total units")
	}

	// regra: mínimo 1 dia, +1 a cada 5 unidades
	deliveryDays := int32(1 + (totalUnits-1)/5)

	return deliveryDays, nil
}

