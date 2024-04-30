package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"log"

	"grpc/internal/pkg/db"
	"grpc/internal/pkg/repository"
)

type InMemoryCache interface {
	GetPickupPoints(id int64) (repository.PickupPoint, error)
	SetPickupPoints(id int64, pickupPoint repository.PickupPoint) error
}

type PickupPointRepo struct {
	db      db.DBops
	IMCache InMemoryCache
}

func NewPickupPoints(database db.DBops) *PickupPointRepo {
	return &PickupPointRepo{
		db: database,
	}
}

func (r *PickupPointRepo) Add(ctx context.Context, pickup_point *repository.PickupPoint) (int64, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var id int64
	err = tx.QueryRow(ctx, `INSERT INTO pickup_points(name, address, phone_number) VALUES ($1,$2,$3) RETURNING id;`, pickup_point.Name, pickup_point.Address, pickup_point.PhoneNumber).Scan(&id)
	if err != nil {
		return 0, err
	}

	err = r.db.Commit(ctx, tx)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PickupPointRepo) GetByID(ctx context.Context, id int64) (*repository.PickupPoint, error) {
	var a repository.PickupPoint
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "SELECT id, name, address, phone_number FROM pickup_points WHERE id=$1", id).Scan(&a.ID, &a.Name, &a.Address, &a.PhoneNumber)
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

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `UPDATE pickup_points SET name=$1, address=$2, phone_number=$3 WHERE id=$4 RETURNING id;`, pickup_point.Name, pickup_point.Address, pickup_point.PhoneNumber, id).Scan(&id)
	if err != nil {
		return err
	}

	err = r.db.Commit(ctx, tx)
	if err != nil {
		return err
	}

	if err := r.IMCache.SetPickupPoints(id, *pickup_point); err != nil {
		log.Println(err)
	}

	return nil
}

func (r *PickupPointRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "DELETE FROM pickup_points where id=$1", id)
	if err != nil {
		return err
	}

	err = r.db.Commit(ctx, tx)
	if err != nil {
		return err
	}
	return nil
}

func (r *PickupPointRepo) List(ctx context.Context) (*[]repository.PickupPoint, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, "SELECT id, name, address, phone_number FROM pickup_points")
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	defer rows.Close()

	var a []repository.PickupPoint
	for rows.Next() {
		var pickupPoint repository.PickupPoint
		if err := rows.Scan(&pickupPoint.ID, &pickupPoint.Name, &pickupPoint.Address, &pickupPoint.PhoneNumber); err != nil {
			return nil, err
		}
		a = append(a, pickupPoint)
	}
	return &a, nil
}
