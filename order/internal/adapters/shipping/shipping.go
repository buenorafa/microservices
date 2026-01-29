package shipping_adapter

import (
	"context"
	"log"
	"time"

	"github.com/buenorafa/microservices-proto/golang/shipping"
	"github.com/buenorafa/microservices/order/internal/application/core/domain"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	shipping shipping.ShippingClient // comes from the generated code by the protobuf compiler
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts, 
	grpc.WithUnaryInterceptor( 
		grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
			grpc_retry.WithMax(5),
			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(1*time.Second)),
		),
	),
	)

	
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(shippingServiceUrl, opts...)
	if err != nil {
		return nil, err
	}

	client := shipping.NewShippingClient(conn) // initialize the stub
	return &Adapter{shipping: client}, nil
}

func (a Adapter) RequestShipping(order *domain.Order) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := a.shipping.Create(ctx, &shipping.CreateShippingRequest{
		OrderId: order.ID,
		Items: order.ToShippingItems(),
	})

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return int32(resp.DeliveryDays), nil
}