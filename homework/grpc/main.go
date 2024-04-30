package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"

	"grpc/internal/app/service"
	"grpc/internal/pkg/db"
	"grpc/internal/pkg/jaeger"
	"grpc/internal/pkg/pb"
	"grpc/internal/pkg/repository/postgresql"
)

func init() {
	// Register standard server metrics and customized metrics to registry.
	service.Reg.MustRegister(service.AddedPointsMetric, service.DeletedPointsMetric, service.InternalErrorCountMetric, service.ClientErrorCountMetric)

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	database, err := db.NewDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()
	srv := service.Service{Repo: postgresql.NewPickupPoints(database)}

	shutdown, err := jaeger.InitProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	srv.Tracer = otel.Tracer("test-tracer")

	// Listen an actual port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9093))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	// Create some standard server metrics.
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	service.Reg.MustRegister(grpcMetrics)

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

	go http.ListenAndServe(":9091", promhttp.HandlerFor(service.Reg, promhttp.HandlerOpts{EnableOpenMetrics: true}))

	// Start your gRPC server.
	log.Fatal(grpcServer.Serve(lis))
}
