package main

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"github.com/mamalmaleki/go_movie/gen"
	"github.com/mamalmaleki/go_movie/metadata/internal/controller/metadata"
	"github.com/mamalmaleki/go_movie/pkg/tracing"
	"github.com/uber-go/tally"
	"github.com/uber-go/tally/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"gopkg.in/yaml.v3"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"

	//"github.com/mamalmaleki/go_movie/metadata/internal/repository/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"

	//httpHandler "github.com/mamalmaleki/go_movie/metadata/internal/handler/http"
	grpcHandler "github.com/mamalmaleki/go_movie/metadata/internal/handler/grpc"
	"github.com/mamalmaleki/go_movie/metadata/internal/repository/mysql"
	"github.com/mamalmaleki/go_movie/pkg/discovery"
	"github.com/mamalmaleki/go_movie/pkg/discovery/consul"
	"go.uber.org/zap"
	"log"
	"time"
)

const serviceName = "metadata"

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("Started the service", zap.String("serviceName", serviceName))

	simulateCPULoad := flag.Bool("simulatecpuload", false, "simulate CPU load for profiling")
	flag.Parse()
	if *simulateCPULoad {
		go heavyOperation()
	}

	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			logger.Fatal("Failed to start profiler handler", zap.Error(err))
		}
	}()

	filename := os.Getenv("CONFIG_FILE")
	if filename == "" {
		var err error
		filename, err = filepath.Abs("../metadata/configs/base.yaml")
		if err != nil {
			panic(err)
		}
	}
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	port := cfg.API.Port
	//flag.IntVar(&port, "port", 8081, "API handler port")
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
	repo, err := mysql.New()
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %v", err)
	}
	ctrl := metadata.New(repo)
	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	reflection.Register(srv)
	gen.RegisterMetadataServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
	//h := httpHandler.New(ctrl)
	//http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	//if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
	//	panic(err)
	//}
}

func heavyOperation() {
	time.Sleep(5 * time.Second)
	token := make([]byte, 8192)
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	random.Read(token)
	md5.New().Write(token)
	log.Println("Heavy")
}
