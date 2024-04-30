package service

import (
	"go.opentelemetry.io/otel/trace"

	"grpc/internal/pkg/pb"
	"grpc/internal/pkg/repository"
)

type Service struct {
	Tracer trace.Tracer
	Repo   repository.PickupPointRepo
	pb.UnimplementedPickupPointsServer
}
