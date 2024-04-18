package interactive

import (
	"testing"

	"github.com/golang/mock/gomock"

	mock_interactive_storage "homework/internal/interactive_storage/mocks"
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
