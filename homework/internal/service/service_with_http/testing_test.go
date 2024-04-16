package service_with_http

import (
	"github.com/golang/mock/gomock"
	mock_repository "homework/internal/pkg/repository/mocks"
	"testing"
)

type pickupPointsRepoFixtures struct {
	ctrl             *gomock.Controller
	srv              Server1
	mockPickupPoints *mock_repository.MockPickupPointRepo
}

func setUp(t *testing.T) pickupPointsRepoFixtures {
	ctrl := gomock.NewController(t)
	mockPickupPoints := mock_repository.NewMockPickupPointRepo(ctrl)
	srv := Server1{mockPickupPoints}
	return pickupPointsRepoFixtures{
		ctrl:             ctrl,
		mockPickupPoints: mockPickupPoints,
		srv:              srv,
	}
}

func (p *pickupPointsRepoFixtures) tearDown() {
	p.ctrl.Finish()
}
