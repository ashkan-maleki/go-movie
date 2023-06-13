package testutil

import (
	"github.com/mamalmaleki/go-movie/gen"
	"github.com/mamalmaleki/go-movie/rating/internal/controller/rating"
	grpcHandler "github.com/mamalmaleki/go-movie/rating/internal/handler/grpc"
	"github.com/mamalmaleki/go-movie/rating/internal/repository/memory"
)

func NewTestRatingGRPCServer() gen.RatingServiceServer {
	r := memory.New()
	ctrl := rating.New(r, nil)
	return grpcHandler.New(ctrl)
}
