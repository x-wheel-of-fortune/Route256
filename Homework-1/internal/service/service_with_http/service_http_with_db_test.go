package service_with_http

import (
	"Homework-1/internal/pkg/repository"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func Test_validateCreate(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		unm = addPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		err := s.srv.validateCreate(ctx, unm)
		assert.Equal(t, nil, err)
	})
}

func Test_create(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		unm = addPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}

		point = &repository.PickupPoint{
			ID:          0,
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		id = int64(1)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPickupPoints.EXPECT().Add(gomock.Any(), point).Return(id, nil)
		result, status, err := s.srv.create(ctx, unm)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "{\"ID\":1,\"name\":\"Name\",\"address\":\"Address\",\"phone_number\":\"PhoneNumber\"}", string(result))
	})
}

func Test_GetByID(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPickupPoints.EXPECT().GetByID(gomock.Any(), id).Return(&repository.PickupPoint{
			ID:          1,
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}, nil)
		result, status := s.srv.get(ctx, id)
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "{\"ID\":1,\"name\":\"Name\",\"address\":\"Address\",\"phone_number\":\"PhoneNumber\"}", string(result))
	})
}
