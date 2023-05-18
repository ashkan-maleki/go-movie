package main

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mamalmaleki/go_movie/gen"
	metadataTest "github.com/mamalmaleki/go_movie/metadata/pkg/testutil"
	movieTest "github.com/mamalmaleki/go_movie/movie/pkg/testutil"
	"github.com/mamalmaleki/go_movie/pkg/discovery"
	"github.com/mamalmaleki/go_movie/pkg/discovery/memory"
	ratingTest "github.com/mamalmaleki/go_movie/rating/pkg/testutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
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

	log.Println("Getting movie details via movie service")
	wantMovieDetails := &gen.MovieDetails{
		Metadata: m,
	}
	getMovieDetailsResponse, err := movieClient.GetMovieDetails(ctx,
		&gen.GetMovieDetailsRequest{MovieId: m.Id})
	if err != nil {
		log.Fatalf("get movie details: %v", err)
	}
	if diff := cmp.Diff(getMovieDetailsResponse.MovieDetails, wantMovieDetails,
		cmpopts.IgnoreUnexported(gen.MovieDetails{}, gen.Metadata{})); diff != "" {
		log.Fatalf("get metadata after put mismatch: %v", diff)
	}

	log.Println("Saving first rating via rating service")
	const userID = "user0"
	const recordTypeMovie = "movie"
	firstRating := int32(5)
	if _, err := ratingClient.PutRating(ctx, &gen.PutRatingRequest{
		UserId:      userID,
		RecordId:    m.Id,
		RecordType:  recordTypeMovie,
		RatingValue: firstRating,
	}); err != nil {
		log.Fatalf("put rating: %v", err)
	}

	log.Println("Retrieving initial aggregated rating via rating service")
	getAggregatedRatingResponse, err := ratingClient.GetAggregatedRating(ctx,
		&gen.GetAggregatedRatingRequest{
			RecordId:   m.Id,
			RecordType: recordTypeMovie,
		})
	if err != nil {
		log.Fatalf("get aggregated rating: %v", err)
	}
	if got, want := getAggregatedRatingResponse.RatingValue, float64(5); got != want {
		log.Fatalf("rating mismatch: got %v want %v", got, want)
	}

	log.Println("Saving second rating via rating service")
	secondRating := int32(1)
	if _, err = ratingClient.PutRating(ctx, &gen.PutRatingRequest{
		UserId:      userID,
		RecordId:    m.Id,
		RecordType:  recordTypeMovie,
		RatingValue: secondRating,
	}); err != nil {
		log.Fatalf("put rating: %v", err)
	}

	log.Println("Getting new aggregated rating via rating service")
	getAggregatedRatingResponse, err = ratingClient.GetAggregatedRating(ctx,
		&gen.GetAggregatedRatingRequest{
			RecordId:   m.Id,
			RecordType: recordTypeMovie,
		})
	if err != nil {
		log.Fatalf("get aggregated rating: %v", err)
	}
	wantRating := float64((firstRating + secondRating) / 2)
	if got, want := getAggregatedRatingResponse.RatingValue, wantRating; got != want {
		log.Fatalf("rating mismatch: got %v want %v", got, want)
	}

	log.Println("Getting updated movie details via movie service")
	getMovieDetailsResponse, err = movieClient.GetMovieDetails(ctx,
		&gen.GetMovieDetailsRequest{MovieId: m.Id})
	if err != nil {
		log.Fatalf("get movie details: %v", err)
	}
	wantMovieDetails.Rating = float32(wantRating)
	if diff := cmp.Diff(getMovieDetailsResponse.MovieDetails, wantMovieDetails,
		cmpopts.IgnoreUnexported(gen.MovieDetails{}, gen.Metadata{})); diff != "" {
		log.Fatalf("get movie details after update mismatch: %v", err)
	}

	log.Println("Integration test execution successful")
}

func startMetadataService(ctx context.Context, registry *memory.Registry) *grpc.Server {
	log.Println("Starting metadata service on " + metadataServiceAddr)
	h := metadataTest.NewTestMovieGRPCServer()
	l, err := net.Listen("tcp", metadataServiceAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, h)
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}
	}()
	id := discovery.GenerateInstanceID(metadataServiceName)
	if err := registry.Register(ctx, id, metadataServiceName, metadataServiceAddr); err != nil {
		panic(err)
	}
	return srv
}

func startRatingService(ctx context.Context, registry *memory.Registry) *grpc.Server {
	log.Println("Starting rating service on " + ratingServiceAddr)
	h := ratingTest.NewTestRatingGRPCServer()
	l, err := net.Listen("tcp", ratingServiceAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}
	}()
	id := discovery.GenerateInstanceID(ratingServiceName)
	if err := registry.Register(ctx, id, ratingServiceName, ratingServiceAddr); err != nil {
		panic(err)
	}
	return srv
}

func startMovieService(ctx context.Context, registry *memory.Registry) *grpc.Server {
	log.Println("Starting movie service on " + movieServiceAddr)
	h := movieTest.NewTestMovieGRPCServer(registry)
	l, err := net.Listen("tcp", movieServiceAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMovieServiceServer(srv, h)
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}
	}()
	id := discovery.GenerateInstanceID(movieServiceName)
	if err := registry.Register(ctx, id, movieServiceName, movieServiceAddr); err != nil {
		panic(err)
	}
	return srv
}
