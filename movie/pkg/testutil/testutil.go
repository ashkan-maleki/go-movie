package testutil

import (
	"github.com/mamalmaleki/go_movie/gen"
	"github.com/mamalmaleki/go_movie/movie/internal/controller/movie"
	metadataGatewayPkg "github.com/mamalmaleki/go_movie/movie/internal/gateway/metadata/grpc"
	ratingGatewayPkg "github.com/mamalmaleki/go_movie/movie/internal/gateway/rating/grpc"
	grpcHandler "github.com/mamalmaleki/go_movie/movie/internal/handler/grpc"
	"github.com/mamalmaleki/go_movie/pkg/discovery"
)

func NewTestMovieGRPCServer(registry discovery.Registry) gen.MovieServiceServer {
	metadataGateway := metadataGatewayPkg.New(registry)
	ratingGateway := ratingGatewayPkg.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	return grpcHandler.New(ctrl)
}
