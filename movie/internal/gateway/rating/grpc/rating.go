package grpc

import (
	"context"
	"github.com/mamalmaleki/go-movie/gen"
	"github.com/mamalmaleki/go-movie/internal/grpcutil"
	"github.com/mamalmaleki/go-movie/movie/internal/gateway"
	"github.com/mamalmaleki/go-movie/pkg/discovery"
	"github.com/mamalmaleki/go-movie/rating/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Gateway defines an gRPC gateway for a rating service.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a rating service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// GetAggregatedRating returns the aggregated rating for a record or ErrNotFound
// if there are no ratings for it.
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID,
	recordType model.RecordType) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	client := gen.NewRatingServiceClient(conn)
	resp, err := client.GetAggregatedRating(ctx,
		&gen.GetAggregatedRatingRequest{RecordId: string(recordID),
			RecordType: string(recordType)})
	stat, ok := status.FromError(err)
	if ok && stat.Code() == codes.NotFound {
		return 0, gateway.ErrNotFound
	} else if err != nil {
		return 0, err
	}
	return resp.RatingValue, nil
}

// PutRating writes a rating.
func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType,
	rating *model.Rating) error {

	return nil
}
