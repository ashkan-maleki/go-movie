package grpc

import (
	"context"
	"github.com/mamalmaleki/go_movie/gen"
	"github.com/mamalmaleki/go_movie/internal/grpcutil"
	"github.com/mamalmaleki/go_movie/metadata/pkg/model"
	"github.com/mamalmaleki/go_movie/pkg/discovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Gateway defines a movie metadata gRPC gateway.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for movie metadata service
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
	if err != nil {
		return nil, err
	}
	return model.MetadataFromProto(resp.Metadata), nil
}

func shouldRetry(err error) bool {
	e, ok := status.FromError(err)
	if !ok {
		return false
	}
	return e.Code() == codes.DeadlineExceeded ||
		e.Code() == codes.ResourceExhausted ||
		e.Code() == codes.Unavailable
}
