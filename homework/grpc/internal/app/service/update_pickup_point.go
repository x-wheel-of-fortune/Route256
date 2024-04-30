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

func (s *Service) UpdatePickupPoint(ctx context.Context, req *pb.PickupPointRequest) (*pb.PickupPointResponse, error) {
	// work begins
	span := trace.SpanFromContext(ctx)
	defer span.End()

	p := req.GetPickupPoint()
	pickupPoint := &repository.PickupPoint{
		ID:          int(p.GetId()),
		Name:        p.GetName(),
		Address:     p.GetAddress(),
		PhoneNumber: p.GetPhoneNumber(),
	}

	err := s.validateUpdate(ctx, pickupPoint)
	if err != nil {
		ClientErrorCountMetric.Add(1)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = s.Repo.Update(ctx, int64(pickupPoint.ID), pickupPoint)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			ClientErrorCountMetric.Add(1)
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		InternalErrorCountMetric.Add(1)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &pb.PickupPointResponse{
		Id:          int64(pickupPoint.ID),
		Name:        pickupPoint.Name,
		Address:     pickupPoint.Address,
		PhoneNumber: pickupPoint.PhoneNumber,
	}

	return resp, nil
}

func (s *Service) validateUpdate(ctx context.Context, pickupPoint *repository.PickupPoint) error {
	if pickupPoint.ID == 0 {
		return errors.New("ID field is empty")
	}
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
