package service

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc/internal/pkg/pb"
)

func (s *Service) DeletePickupPoint(ctx context.Context, req *pb.IdRequest) (*pb.Empty, error) {
	// work begins
	span := trace.SpanFromContext(ctx)
	defer span.End()
	if req.GetId() == 0 {
		ClientErrorCountMetric.Add(1)
		return nil, status.Errorf(codes.InvalidArgument, "id not specified")
	}
	err := s.Repo.Delete(ctx, req.GetId())
	if err != nil {
		InternalErrorCountMetric.Add(1)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.Empty{}, nil
}
