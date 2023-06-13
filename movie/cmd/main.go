package main

import (
	"context"
	"fmt"
	"github.com/mamalmaleki/go-movie/gen"
	"github.com/mamalmaleki/go-movie/movie/internal/controller/movie"
	grpcHandler "github.com/mamalmaleki/go-movie/movie/internal/handler/grpc"
	"github.com/mamalmaleki/go-movie/pkg/tracing"
	"github.com/uber-go/tally"
	"github.com/uber-go/tally/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"net"
	"net/http"
	"os"
	"path/filepath"

	//metadataGatewayPkg "github.com/mamalmaleki/go-movie/movie/internal/gateway/metadata/http"
	metadataGatewayPkg "github.com/mamalmaleki/go-movie/movie/internal/gateway/metadata/grpc"
	ratingGatewayPkg "github.com/mamalmaleki/go-movie/movie/internal/gateway/rating/grpc"

	"github.com/mamalmaleki/go-movie/pkg/discovery"
	"github.com/mamalmaleki/go-movie/pkg/discovery/consul"
	"log"
	"time"
)

const serviceName = "movie"

func main() {
	logger, _ := zap.NewProduction()
	logger.Info("Started the service", zap.String("serviceName", serviceName))

	filename := os.Getenv("CONFIG_FILE")
	if filename == "" {
		var err error
		filename, err = filepath.Abs("../movie/configs/base.yaml")
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
		Tags:           map[string]string{"service": serviceName},
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
		"service": serviceName,
	}).Counter("service_started")
	counter.Inc(1)

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
	//const limit = 100
	//const burst = 100
	//l := newLimiter(limit, burst)
	srv := grpc.NewServer(
		//grpc.UnaryInterceptor(ratelimit.UnaryServerInterceptor(l)),
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)
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
