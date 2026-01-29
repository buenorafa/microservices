package main

import (
	"log"
	"net"
	"os"

	shippingpb "github.com/buenorafa/microservices-proto/golang/shipping"
	grpcadapter "github.com/buenorafa/microservices/shipping/internal/adapters/grpc"
	"github.com/buenorafa/microservices/shipping/internal/application/core/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("APPLICATION_PORT")
	if port == "" {
		log.Fatal("APPLICATION_PORT environment variable is missing")
	}

	app := api.NewApplication()

	server := grpc.NewServer()
	handler := grpcadapter.NewAdapter(app) // ou NewShippingServer(app), conforme seu nome

	shippingpb.RegisterShippingServer(server, handler)
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("starting shipping service on port %s ...", port)
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}