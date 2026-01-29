package main

import (
	"log"

	"github.com/buenorafa/microservices/order/config"
	"github.com/buenorafa/microservices/order/internal/adapters/db"
	"github.com/buenorafa/microservices/order/internal/adapters/grpc"
	payment_adapter "github.com/buenorafa/microservices/order/internal/adapters/payment"
	shipping_adapter "github.com/buenorafa/microservices/order/internal/adapters/shipping"
	"github.com/buenorafa/microservices/order/internal/application/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("failed to connect to database. Error: %v", err)
	}

	paymentAdapter, err := payment_adapter.NewAdapter(config.GetPaymentServiceURL())
	if err != nil {
		log.Fatalf("failed to initialize payment stub. Error: %v", err)
	}

	shippingAdapter, err := shipping_adapter.NewAdapter(config.GetShippingServiceURL())
	if err != nil {
		log.Fatalf("failed to initialize payment stub. Error: %v", err)
	}

	application := api.NewApplication(dbAdapter, paymentAdapter, shippingAdapter)

	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}