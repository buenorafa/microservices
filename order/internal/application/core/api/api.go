package api

import (
	"github.com/buenorafa/microservices/order/internal/application/core/domain"
	"github.com/buenorafa/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db: db,
		payment: payment,
	}
}

func (a *Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	
	if order.TotalItems() > 50 {
		return domain.Order{}, status.Errorf(codes.InvalidArgument, "Order over 50 items is not allowed.")
	}

	order.Status = "Pending"

	if err := a.db.Save(&order); err != nil {
		// não existe registro para atualizar status
		return domain.Order{}, err
	}

	if err := a.payment.Charge(&order); err != nil {
		order.Status = "Canceled"
		_ = a.db.UpdateStatus(&order)
		return domain.Order{}, err
	}

	order.Status = "Paid"
	if err := a.db.UpdateStatus(&order); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}