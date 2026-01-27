package payment_adapter

import (
	"context"
	"log"
	"time"

	"github.com/buenorafa/microservices-proto/golang/payment"
	"github.com/buenorafa/microservices/order/internal/application/core/domain"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	payment payment.PaymentClient // comes from the generated code by the protobuf compiler
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
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

	conn, err := grpc.Dial(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}

	client := payment.NewPaymentClient(conn) // initialize the stub
	return &Adapter{payment: client}, nil
}

func (a Adapter) Charge(order *domain.Order) error {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	_, err := a.payment.Create(ctx, &payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})

	if err != nil {
		log.Println(err)
	}

	return err
}