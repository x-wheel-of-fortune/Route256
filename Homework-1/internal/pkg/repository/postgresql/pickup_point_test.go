package postgresql

import (
	mock_database "Homework-1/internal/pkg/db/mocks"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Add(t *testing.T) {
	t.Parallel()
}

func Test_GetByID(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = 0
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockDB := mock_database.NewMockDBops(ctrl)
		repo := NewPickupPoints(mockDB)
		mockDB.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,name,address, phone_number FROM pickup_points where id=$1", gomock.Any()).Return(nil)

		// act
		pickupPoint, err := repo.GetByID(ctx, int64(id))

		// assert
		require.NoError(t, err)
		assert.Equal(t, 0, pickupPoint.ID)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("not found", func(t *testing.T) {
			t.Parallel()

		})
		t.Run("internal error", func(t *testing.T) {
			t.Parallel()

		})
	})
}
