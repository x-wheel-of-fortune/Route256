package service

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc/internal/pkg/pb"
)

func (s *Service) GetPickupPoint(ctx context.Context, req *pb.IdRequest) (*pb.PickupPointResponse, error) {
	// work begins
	span := trace.SpanFromContext(ctx)
	defer span.End()
	if req.GetId() == 0 {
		ClientErrorCountMetric.Add(1)
		return nil, status.Errorf(codes.InvalidArgument, "id not specified")
	}
	point, err := s.Repo.GetByID(ctx, req.GetId())
	if err != nil {
		InternalErrorCountMetric.Add(1)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &pb.PickupPointResponse{
		Id:          int64(point.ID),
		Name:        point.Name,
		Address:     point.Address,
		PhoneNumber: point.PhoneNumber,
	}

	return resp, nil
}
