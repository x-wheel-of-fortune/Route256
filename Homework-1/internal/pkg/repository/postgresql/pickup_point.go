package postgresql

import (
	"Homework-1/internal/pkg/db"
	"Homework-1/internal/pkg/repository"
	"context"
	"database/sql"
	"errors"
)

type PickupPointRepo struct {
	db *db.Database
}

func NewPickupPoints(database *db.Database) *PickupPointRepo {
	return &PickupPointRepo{db: database}
}

func (r *PickupPointRepo) Add(ctx context.Context, pickup_point *repository.PickupPoint) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO pickup_points(name, address, phone_number) VALUES ($1,$2,$3) RETURNING id;`, pickup_point.Name, pickup_point.Address, pickup_point.PhoneNumber).Scan(&id)
	return id, err
}

func (r *PickupPointRepo) GetByID(ctx context.Context, id int64) (*repository.PickupPoint, error) {
	var a repository.PickupPoint
	err := r.db.Get(ctx, &a, "SELECT id,name,address, phone_number FROM pickup_points where id=$1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *PickupPointRepo) Delete(ctx context.Context) error {

	return nil
}
