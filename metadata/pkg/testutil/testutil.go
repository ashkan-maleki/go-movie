package testutil

import (
	"github.com/mamalmaleki/go-movie/gen"
	"github.com/mamalmaleki/go-movie/metadata/internal/controller/metadata"
	grpcHandler "github.com/mamalmaleki/go-movie/metadata/internal/handler/grpc"
	"github.com/mamalmaleki/go-movie/metadata/internal/repository/memory"
)

func NewTestMovieGRPCServer() gen.MetadataServiceServer {
	r := memory.New()
	ctrl := metadata.New(r)
	return grpcHandler.New(ctrl)
}
