package interactive

import (
	mock_interactive_storage "Homework-1/internal/interactive_storage/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

type interactiveStorageFixtures struct {
	ctrl                   *gomock.Controller
	mockInteractiveStorage *mock_interactive_storage.MockInteractiveStorage
}

func setUp(t *testing.T) interactiveStorageFixtures {
	ctrl := gomock.NewController(t)
	mockStorage := mock_interactive_storage.NewMockInteractiveStorage(ctrl)
	return interactiveStorageFixtures{
		ctrl:                   ctrl,
		mockInteractiveStorage: mockStorage,
	}
}

func (s *interactiveStorageFixtures) tearDown() {
	s.ctrl.Finish()
}
