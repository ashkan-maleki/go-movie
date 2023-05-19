package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"github.com/mamalmaleki/go_movie/gen"
	"github.com/mamalmaleki/go_movie/movie/internal/controller/movie"
	grpcHandler "github.com/mamalmaleki/go_movie/movie/internal/handler/grpc"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"net"
	"os"

	//metadataGatewayPkg "github.com/mamalmaleki/go_movie/movie/internal/gateway/metadata/http"
	metadataGatewayPkg "github.com/mamalmaleki/go_movie/movie/internal/gateway/metadata/grpc"
	ratingGatewayPkg "github.com/mamalmaleki/go_movie/movie/internal/gateway/rating/grpc"

	"github.com/mamalmaleki/go_movie/pkg/discovery"
	"github.com/mamalmaleki/go_movie/pkg/discovery/consul"
	"log"
	"time"
)

const serviceName = "movie"

func main() {
	log.Println("Starting the movie gateway service")
	filename := os.Getenv("CONFIG_FILE")
	if filename == "" {
		filename = "./movie/configs/base.yaml"
	}
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	port := cfg.API.Port
	log.Printf("Starting the movie metadata service on port %d", port)
	serviceDiscoverUrl := os.Getenv("SERVICE_DISCOVERY_URL")
	if serviceDiscoverUrl == "" {
		serviceDiscoverUrl = "localhost:8500"
	}
	registry, err := consul.NewRegistry(serviceDiscoverUrl)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName,
		fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	metadataGateway := metadataGatewayPkg.New(registry)
	ratingGateway := ratingGatewayPkg.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)

	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	const limit = 100
	const burst = 100
	l := newLimiter(limit, burst)
	srv := grpc.NewServer(grpc.UnaryInterceptor(ratelimit.UnaryServerInterceptor(l)))
	reflection.Register(srv)
	gen.RegisterMovieServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}

	//h := httpHandler.New(ctrl)
	//http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	//if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
	//	panic(err)
	//}
}

type limiter struct {
	l *rate.Limiter
}

func newLimiter(limit int, burst int) *limiter {
	return &limiter{rate.NewLimiter(rate.Limit(limit), burst)}
}

func (l *limiter) Limit() bool {
	return l.l.Allow()
}
