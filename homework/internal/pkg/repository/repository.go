//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository

package repository

import "context"

type PickupPointRepo interface {
	Add(ctx context.Context, pickup_point *PickupPoint) (int64, error)
	GetByID(ctx context.Context, id int64) (*PickupPoint, error)
	Update(ctx context.Context, id int64, pickup_point *PickupPoint) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) (*[]PickupPoint, error)
}
