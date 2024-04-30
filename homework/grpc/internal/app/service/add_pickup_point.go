package service

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc/internal/pkg/pb"
	"grpc/internal/pkg/repository"
)

func (s *Service) AddPickupPoint(ctx context.Context, req *pb.PickupPointRequest) (*pb.PickupPointResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	p := req.GetPickupPoint()
	pickupPoint := &repository.PickupPoint{
		Name:        p.GetName(),
		Address:     p.GetAddress(),
		PhoneNumber: p.GetPhoneNumber(),
	}
	err := s.validateAdd(ctx, pickupPoint)
	if err != nil {
		ClientErrorCountMetric.Add(1)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := s.Repo.Add(ctx, pickupPoint)
	if err != nil {
		InternalErrorCountMetric.Add(1)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &pb.PickupPointResponse{
		Id:          id,
		Name:        pickupPoint.Name,
		Address:     pickupPoint.Address,
		PhoneNumber: pickupPoint.PhoneNumber,
	}

	AddedPointsMetric.Add(1)
	return resp, nil
}

func (s *Service) validateAdd(ctx context.Context, pickupPoint *repository.PickupPoint) error {
	if pickupPoint.Name == "" {
		return errors.New("Name field is empty")
	}
	if pickupPoint.Address == "" {
		return errors.New("Address field is empty")
	}
	if pickupPoint.PhoneNumber == "" {
		return errors.New("PhoneNumber field is empty")
	}
	return nil
}
