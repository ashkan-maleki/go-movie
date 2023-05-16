package main

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mamalmaleki/go_movie/gen"
	"github.com/mamalmaleki/go_movie/pkg/discovery/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const (
	metadataServiceName = "metadata"
	ratingServiceName   = "rating"
	movieServiceName    = "movie"

	metadataServiceAddr = "localhost:8081"
	ratingServiceAddr   = "localhost:8082"
	movieServiceAddr    = "localhost:8083"
)

func main() {
	log.Println("Starting the integration test")

	ctx := context.Background()
	registry := memory.NewRegistry()

	log.Println("Setting up service handlers and clients")

	metadataService := startMetadataService(ctx, registry)
	defer metadataService.GracefulStop()
	ratingService := startRatingService(ctx, registry)
	defer ratingService.GracefulStop()
	movieService := startMovieService(ctx, registry)
	defer movieService.GracefulStop()

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	metadataConnection, err := grpc.Dial(metadataServiceAddr, opts)
	if err != nil {
		panic(err)
	}
	defer metadataConnection.Close()
	metadataClient := gen.NewMetadataServiceClient(metadataConnection)

	ratingConnection, err := grpc.Dial(ratingServiceAddr, opts)
	if err != nil {
		panic(err)
	}
	defer ratingConnection.Close()
	ratingClient := gen.NewRatingServiceClient(ratingConnection)

	movieConnection, err := grpc.Dial(movieServiceAddr, opts)
	if err != nil {
		panic(err)
	}
	defer movieConnection.Close()
	movieClient := gen.NewMovieServiceClient(movieConnection)

	log.Println("Saving test metadata via metadata service")
	m := &gen.Metadata{
		Id:          "the-movie",
		Title:       "The Movie",
		Description: "The Movie, the one and only",
		Director:    "Mr. D",
	}
	if _, err := metadataClient.PutMetadata(ctx, &gen.PutMetadataRequest{Metadata: m}); err != nil {
		log.Fatalf("put metadata: %v", err)
	}

	log.Println("Retrieving test metadata via metadata service")
	getMetadataResponse, err := metadataClient.GetMetadata(ctx,
		&gen.GetMetadataRequest{MovieId: m.Id})
	if err != nil {
		log.Fatalf("get metadata: %v", err)
	}
	if diff := cmp.Diff(getMetadataResponse.Metadata, m,
		cmpopts.IgnoreUnexported(gen.Metadata{})); diff != "" {
		log.Fatalf("get metadata after put mismatch: %v", diff)
	}
}

func startMetadataService(ctx context.Context, registry *memory.Registry) *grpc.Server {
	// TODO
	return nil
}

func startRatingService(ctx context.Context, registry *memory.Registry) *grpc.Server {
	// TODO
	return nil
}

func startMovieService(ctx context.Context, registry *memory.Registry) *grpc.Server {
	// TODO
	return nil
}
