package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4"
	"homework/internal/pkg/db"
	"homework/internal/pkg/repository"
)

type PickupPointRepo struct {
	db db.DBops
}

func NewPickupPoints(database db.DBops) *PickupPointRepo {
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *PickupPointRepo) Update(ctx context.Context, id int64, pickup_point *repository.PickupPoint) error {
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = r.db.ExecQueryRow(ctx, `UPDATE pickup_points SET name=$1, address=$2, phone_number=$3 WHERE id=$4 RETURNING id;`, pickup_point.Name, pickup_point.Address, pickup_point.PhoneNumber, id).Scan(&id)
	if err != nil {
		return err
	}
	return err
}

func (r *PickupPointRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, "DELETE FROM pickup_points where id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PickupPointRepo) List(ctx context.Context) (*[]repository.PickupPoint, error) {
	var a []repository.PickupPoint
	err := r.db.Select(ctx, &a, "SELECT id, name, address, phone_number FROM pickup_points")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	return &a, nil
}
