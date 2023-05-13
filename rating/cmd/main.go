package main

import (
	"context"
	"fmt"
	"github.com/mamalmaleki/go_movie/gen"
	"github.com/mamalmaleki/go_movie/pkg/discovery"
	"github.com/mamalmaleki/go_movie/pkg/discovery/consul"
	"github.com/mamalmaleki/go_movie/rating/internal/controller/rating"
	grpcHandler "github.com/mamalmaleki/go_movie/rating/internal/handler/grpc"
	"github.com/mamalmaleki/go_movie/rating/internal/ingester/kafka"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"net"
	"os"

	//httpHandler "github.com/mamalmaleki/go_movie/rating/internal/handler/http"
	//"github.com/mamalmaleki/go_movie/rating/internal/repository/memory"
	"github.com/mamalmaleki/go_movie/rating/internal/repository/mysql"
	"log"
	"time"
)

const serviceName = "rating"

func main() {
	log.Println("Starting the movie rating service")
	filename := os.Getenv("CONFIG_FILE")
	if filename == "" {
		filename = "./rating/configs/base.yaml"
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
	//flag.IntVar(&port, "port", 8082, "API handler port")
	//flag.Parse()
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

	kafkaAddress := os.Getenv("KAFKA_ADDRESS")
	if kafkaAddress == "" {
		kafkaAddress = "localhost"
	}
	ingester, err := kafka.NewIngester(kafkaAddress, "my-group", "ratings")
	if err != nil {
		log.Fatalf("failed to create ingester: %v", err)
	}
	repo, err := mysql.New()
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %v", err)
	}
	ctrl := rating.New(repo, ingester)
	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterRatingServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
	//h := httpHandler.New(ctrl)
	//http.Handle("/rating", http.HandlerFunc(h.Handle))
	//if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
	//	panic(err)
	//}
}
