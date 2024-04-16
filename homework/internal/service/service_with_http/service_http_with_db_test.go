package service_with_http

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework/internal/pkg/repository"
	"net/http"
	"testing"
)

func Test_validateCreate(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		unm := AddPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		s := setUp(t)
		defer s.tearDown()
		err := s.srv.validateCreate(ctx, unm)
		assert.Equal(t, nil, err)
	})

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()
		unm := AddPickupPointRequest{
			Name:        "",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		s := setUp(t)
		defer s.tearDown()
		err := s.srv.validateCreate(ctx, unm)
		assert.EqualError(t, err, "Name field is empty")
	})

	t.Run("empty address", func(t *testing.T) {
		t.Parallel()
		unm := AddPickupPointRequest{
			Name:        "Name",
			Address:     "",
			PhoneNumber: "PhoneNumber",
		}
		s := setUp(t)
		defer s.tearDown()
		err := s.srv.validateCreate(ctx, unm)
		assert.EqualError(t, err, "Address field is empty")
	})

	t.Run("empty phoneNumber", func(t *testing.T) {
		t.Parallel()
		unm := AddPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "",
		}
		s := setUp(t)
		defer s.tearDown()
		err := s.srv.validateCreate(ctx, unm)
		assert.EqualError(t, err, "PhoneNumber field is empty")
	})
}

func Test_create(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		unm = AddPickupPointRequest{
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
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "{\"ID\":1,\"name\":\"Name\",\"address\":\"Address\",\"phone_number\":\"PhoneNumber\"}", string(result))
	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPickupPoints.EXPECT().Add(gomock.Any(), point).Return(id, errors.New("internal error"))
		result, status, err := s.srv.create(ctx, unm)
		assert.EqualError(t, err, "internal error")
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Equal(t, "", string(result))
	})
}

func Test_validateUpdate(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		unm := UpdatePickupPointRequest{
			ID:          1,
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		s := setUp(t)
		defer s.tearDown()
		err := s.srv.validateUpdate(ctx, unm)
		require.NoError(t, err)
	})

	t.Run("empty id", func(t *testing.T) {
		t.Parallel()
		unm := UpdatePickupPointRequest{
			ID:          0,
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		s := setUp(t)
		defer s.tearDown()
		err := s.srv.validateUpdate(ctx, unm)
		assert.EqualError(t, err, "ID field is empty")
	})

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()
		unm := UpdatePickupPointRequest{
			ID:          1,
			Name:        "",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		s := setUp(t)
		defer s.tearDown()
		err := s.srv.validateUpdate(ctx, unm)
		assert.EqualError(t, err, "Name field is empty")
	})

	t.Run("empty address", func(t *testing.T) {
		t.Parallel()
		unm := UpdatePickupPointRequest{
			ID:          1,
			Name:        "Name",
			Address:     "",
			PhoneNumber: "PhoneNumber",
		}
		s := setUp(t)
		defer s.tearDown()
		err := s.srv.validateUpdate(ctx, unm)
		assert.EqualError(t, err, "Address field is empty")
	})

	t.Run("empty phoneNumber", func(t *testing.T) {
		t.Parallel()
		unm := UpdatePickupPointRequest{
			ID:          1,
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "",
		}
		s := setUp(t)
		defer s.tearDown()
		err := s.srv.validateUpdate(ctx, unm)
		assert.EqualError(t, err, "PhoneNumber field is empty")
	})
}

func Test_update(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		unm = UpdatePickupPointRequest{
			ID:          1,
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
		s.mockPickupPoints.EXPECT().Update(gomock.Any(), id, point).Return(nil)
		result, status, err := s.srv.update(ctx, unm)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "{\"ID\":1,\"name\":\"Name\",\"address\":\"Address\",\"phone_number\":\"PhoneNumber\"}", string(result))
	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPickupPoints.EXPECT().Update(gomock.Any(), id, point).Return(errors.New("internal error"))
		result, status, err := s.srv.update(ctx, unm)
		assert.EqualError(t, err, "internal error")
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Equal(t, "", string(result))
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPickupPoints.EXPECT().Update(gomock.Any(), id, point).Return(repository.ErrObjectNotFound)
		result, status, err := s.srv.update(ctx, unm)
		assert.EqualError(t, err, "Could not find object with id =1")
		assert.Equal(t, http.StatusNotFound, status)
		assert.Equal(t, "", string(result))
	})
}

func Test_get(t *testing.T) {
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
		result, status, err := s.srv.get(ctx, id)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "{\"ID\":1,\"name\":\"Name\",\"address\":\"Address\",\"phone_number\":\"PhoneNumber\"}", string(result))
	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPickupPoints.EXPECT().GetByID(gomock.Any(), id).Return(&repository.PickupPoint{}, errors.New("internal error"))
		result, status, err := s.srv.get(ctx, id)
		assert.EqualError(t, err, "internal error")
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Equal(t, "", string(result))
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPickupPoints.EXPECT().GetByID(gomock.Any(), id).Return(&repository.PickupPoint{}, repository.ErrObjectNotFound)
		result, status, err := s.srv.get(ctx, id)
		assert.EqualError(t, err, "Could not find object with id =1")
		assert.Equal(t, http.StatusNotFound, status)
		assert.Equal(t, "", string(result))
	})
}

func Test_delete(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPickupPoints.EXPECT().Delete(gomock.Any(), id).Return(nil)
		status, err := s.srv.delete(ctx, id)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)
	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPickupPoints.EXPECT().Delete(gomock.Any(), id).Return(errors.New("internal error"))
		status, err := s.srv.delete(ctx, id)
		assert.EqualError(t, err, "internal error")
		assert.Equal(t, http.StatusInternalServerError, status)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPickupPoints.EXPECT().Delete(gomock.Any(), id).Return(repository.ErrObjectNotFound)
		status, err := s.srv.delete(ctx, id)
		assert.EqualError(t, err, "Could not find object with id =1")
		assert.Equal(t, http.StatusNotFound, status)
	})
}

func Test_list(t *testing.T) {
	t.Parallel()
	var (
		ctx    = context.Background()
		point1 = repository.PickupPoint{
			ID:          1,
			Name:        "Name1",
			Address:     "Address1",
			PhoneNumber: "PhoneNumber1",
		}
		point2 = repository.PickupPoint{
			ID:          2,
			Name:        "Name2",
			Address:     "Address2",
			PhoneNumber: "PhoneNumber2",
		}
		point3 = repository.PickupPoint{
			ID:          3,
			Name:        "Name3",
			Address:     "Address3",
			PhoneNumber: "PhoneNumber3",
		}

		points = &[]repository.PickupPoint{point1, point2, point3}
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPickupPoints.EXPECT().List(gomock.Any()).Return(points, nil)
		result, status, err := s.srv.list(ctx)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "[{\"ID\":1,\"name\":\"Name1\",\"address\":\"Address1\",\"phone_number\":\"PhoneNumber1\"},{\"ID\":2,\"name\":\"Name2\",\"address\":\"Address2\",\"phone_number\":\"PhoneNumber2\"},{\"ID\":3,\"name\":\"Name3\",\"address\":\"Address3\",\"phone_number\":\"PhoneNumber3\"}]", string(result))
	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPickupPoints.EXPECT().List(gomock.Any()).Return(&[]repository.PickupPoint{}, errors.New("internal error"))
		result, status, err := s.srv.list(ctx)
		assert.EqualError(t, err, "internal error")
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Equal(t, "", string(result))
	})
}
