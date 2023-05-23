package main

import (
	"context"
	"fmt"
	"github.com/mamalmaleki/go_movie/gen"
	"github.com/mamalmaleki/go_movie/pkg/discovery"
	"github.com/mamalmaleki/go_movie/pkg/discovery/consul"
	"github.com/mamalmaleki/go_movie/pkg/tracing"
	"github.com/mamalmaleki/go_movie/rating/internal/controller/rating"
	grpcHandler "github.com/mamalmaleki/go_movie/rating/internal/handler/grpc"
	"github.com/mamalmaleki/go_movie/rating/internal/ingester/kafka"
	"github.com/uber-go/tally"
	"github.com/uber-go/tally/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	//httpHandler "github.com/mamalmaleki/go_movie/rating/internal/handler/http"
	//"github.com/mamalmaleki/go_movie/rating/internal/repository/memory"
	"github.com/mamalmaleki/go_movie/rating/internal/repository/mysql"
	"log"
	"time"
)

const serviceName = "rating"

func main() {
	logger, _ := zap.NewProduction()
	logger.Info("Started the service", zap.String("serviceName", serviceName))

	filename := os.Getenv("CONFIG_FILE")
	if filename == "" {
		//filename = "../configs/base.yaml"
		//filename = "./rating/configs/base.yaml"
		//filename, _ = os.Getwd()
		var err error
		filename, err = filepath.Abs("../rating/configs/base.yaml")
		if err != nil {
			panic(err)
		}
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp, err := tracing.NewJaegerProvider(cfg.Jaeger.URL, serviceName)
	if err != nil {
		log.Fatal("Failed to initialize Jaeger provider", zap.Error(err))
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal("Failed to shut down Jaeger provider", zap.Error(err))
		}
	}()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	reporter := prometheus.NewReporter(prometheus.Options{})
	scope, closer := tally.NewRootScope(tally.ScopeOptions{
		Tags:           map[string]string{"service": "metadata"},
		CachedReporter: reporter,
	}, 10*time.Second)
	defer closer.Close()
	http.Handle("/metrics", reporter.HTTPHandler())

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d",
			cfg.Prometheus.MetricsPort), nil); err != nil {
			logger.Fatal("Failed to start the metrics handler", zap.Error(err))
		}
	}()

	counter := scope.Tagged(map[string]string{
		"service": "metadata",
	}).Counter("service_started")
	counter.Inc(1)

	serviceDiscoverUrl := os.Getenv("SERVICE_DISCOVERY_URL")
	if serviceDiscoverUrl == "" {
		serviceDiscoverUrl = "localhost:8500"
	}
	registry, err := consul.NewRegistry(serviceDiscoverUrl)
	if err != nil {
		panic(err)
	}

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
	srv := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	reflection.Register(srv)
	gen.RegisterRatingServiceServer(srv, h)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		s := <-sigChan
		cancel()
		log.Printf("Received signal %v, attemting graceful shutdown", s)
		srv.GracefulStop()
		log.Println("Gracefully stopped the gRPC server")
	}()

	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
	wg.Wait()
	//h := httpHandler.New(ctrl)
	//http.Handle("/rating", http.HandlerFunc(h.Handle))
	//if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
	//	panic(err)
	//}
}
