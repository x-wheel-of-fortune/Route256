package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"grpc/internal/pkg/db"
	"grpc/internal/pkg/pb"
	"grpc/internal/pkg/repository"
	"grpc/internal/pkg/repository/postgresql"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

type Service struct {
	tracer trace.Tracer
	Repo   repository.PickupPointRepo
	pb.UnimplementedPickupPointsServer
}

var (
	// Create a metrics registry.
	reg = prometheus.NewRegistry()

	// Create a customized counter metric.
	customizedCounterMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "added_pickup_point_count",
		Help: "Total number of pickup points added.",
	})
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName("test-service"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint("localhost:16686"))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

func init() {
	// Register standard server metrics and customized metrics to registry.
	reg.MustRegister(customizedCounterMetric)

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func (s *Service) AddPickupPoint(ctx context.Context, req *pb.PickupPointRequest) (*pb.PickupPointResponse, error) {
	// work begins
	span := trace.SpanFromContext(ctx)
	defer span.End()
	defer customizedCounterMetric.Add(1)

	pickupPoint := &repository.PickupPoint{
		Name:        req.PickupPoint.Name,
		Address:     req.PickupPoint.Address,
		PhoneNumber: req.PickupPoint.PhoneNumber,
	}

	id, err := s.Repo.Add(ctx, pickupPoint)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			log.Println(err)
		}
		log.Println(err)
	}
	resp := &pb.PickupPointResponse{
		Id:          id,
		Name:        pickupPoint.Name,
		Address:     pickupPoint.Address,
		PhoneNumber: pickupPoint.PhoneNumber,
	}

	return resp, nil
}

func (s *Service) UpdatePickupPoint(ctx context.Context, req *pb.PickupPointRequest) (*pb.PickupPointResponse, error) {
	// work begins
	span := trace.SpanFromContext(ctx)
	defer span.End()

	pickupPoint := &repository.PickupPoint{
		ID:          int(req.PickupPoint.Id),
		Name:        req.PickupPoint.Name,
		Address:     req.PickupPoint.Address,
		PhoneNumber: req.PickupPoint.PhoneNumber,
	}

	err := s.Repo.Update(ctx, req.PickupPoint.Id, pickupPoint)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			log.Println(err)
		}
		log.Println(err)
	}
	resp := &pb.PickupPointResponse{
		Id:          req.PickupPoint.Id,
		Name:        pickupPoint.Name,
		Address:     pickupPoint.Address,
		PhoneNumber: pickupPoint.PhoneNumber,
	}

	return resp, nil
}

func (s *Service) GetPickupPoint(ctx context.Context, req *pb.IdRequest) (*pb.PickupPointResponse, error) {
	// work begins
	span := trace.SpanFromContext(ctx)
	defer span.End()

	point, err := s.Repo.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			log.Println(err)
		}
		log.Println(err)
	}
	resp := &pb.PickupPointResponse{
		Id:          int64(point.ID),
		Name:        point.Name,
		Address:     point.Address,
		PhoneNumber: point.PhoneNumber,
	}

	return resp, nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	database, err := db.NewDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()
	srv := Service{Repo: postgresql.NewPickupPoints(database)}

	shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	srv.tracer = otel.Tracer("test-tracer")

	// Listen an actual port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9093))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	// Create some standard server metrics.
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	reg.MustRegister(grpcMetrics)

	// Create a gRPC Server with gRPC interceptor.
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			grpcMetrics.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			grpcMetrics.StreamServerInterceptor(),
		),
	)

	pb.RegisterPickupPointsServer(grpcServer, &srv)
	grpcMetrics.InitializeMetrics(grpcServer)

	go http.ListenAndServe(":9091", promhttp.HandlerFor(reg, promhttp.HandlerOpts{EnableOpenMetrics: true}))

	// Start your gRPC server.
	log.Fatal(grpcServer.Serve(lis))
}
