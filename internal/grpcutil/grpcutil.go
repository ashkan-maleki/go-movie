package grpcutil

import (
	"context"
	"github.com/mamalmaleki/go-movie/pkg/discovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
)

// ServiceConnection attempts to select a random service instance and returns a gRPC connection to it.
func ServiceConnection(ctx context.Context, serviceName string,
	registry discovery.Registry) (*grpc.ClientConn, error) {
	address, err := registry.ServiceAddresses(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(address[rand.Intn(len(address))],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
}
