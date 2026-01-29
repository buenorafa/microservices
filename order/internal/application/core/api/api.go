package api

import (
	"log"

	"github.com/buenorafa/microservices/order/internal/application/core/domain"
	"github.com/buenorafa/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db ports.DBPort
	payment ports.PaymentPort
	shipping ports.ShippingPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort, shipping ports.ShippingPort) *Application {
	return &Application{
		db: db,
		payment: payment,
		shipping: shipping,
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

	deliveryDays, err := a.shipping.RequestShipping(&order)
	if err != nil {
		log.Println("shipping error: ", err)
	} else {
		log.Println("delivery days: ", deliveryDays)
	}

	return order, nil
}