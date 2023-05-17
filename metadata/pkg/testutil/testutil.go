package testutil

import (
	"github.com/mamalmaleki/go_movie/gen"
	"github.com/mamalmaleki/go_movie/metadata/internal/controller/metadata"
	grpcHandler "github.com/mamalmaleki/go_movie/metadata/internal/handler/grpc"
	"github.com/mamalmaleki/go_movie/metadata/internal/repository/memory"
)

func NewTestMovieGRPCServer() gen.MetadataServiceServer {
	r := memory.New()
	ctrl := metadata.New(r)
	return grpcHandler.New(ctrl)
}
