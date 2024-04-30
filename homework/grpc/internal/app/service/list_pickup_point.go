package service

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc/internal/pkg/pb"
)

func (s *Service) ListPickupPoint(ctx context.Context, req *pb.Empty) (*pb.ListPickupPointResponse, error) {
	// work begins
	span := trace.SpanFromContext(ctx)
	defer span.End()

	points, err := s.Repo.List(ctx)
	if err != nil {
		InternalErrorCountMetric.Add(1)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var pickupPoints []*pb.PickupPoint
	for _, point := range *points {
		pickupPoints = append(pickupPoints, &pb.PickupPoint{
			Id:          int64(point.ID),
			Name:        point.Name,
			Address:     point.Address,
			PhoneNumber: point.PhoneNumber,
		})
	}
	resp := &pb.ListPickupPointResponse{PickupPoints: pickupPoints}
	return resp, nil
}
