package service_with_http

import (
	mock_repository "Homework-1/internal/pkg/repository/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

type pickupPointsRepoFixtures struct {
	ctrl             *gomock.Controller
	srv              server1
	mockPickupPoints *mock_repository.MockPickupPointRepo
}

func setUp(t *testing.T) pickupPointsRepoFixtures {
	ctrl := gomock.NewController(t)
	mockPickupPoints := mock_repository.NewMockPickupPointRepo(ctrl)
	srv := server1{mockPickupPoints}
	return pickupPointsRepoFixtures{
		ctrl:             ctrl,
		mockPickupPoints: mockPickupPoints,
		srv:              srv,
	}
}

func (p *pickupPointsRepoFixtures) tearDown() {
	p.ctrl.Finish()
}
