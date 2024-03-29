package repository

import (
	"errors"
	"time"
)

var ErrObjectNotFound = errors.New("not found")

type Article struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Rating    int64     `db:"rating"`
	CreatedAt time.Time `db:"-"`
}

type PickupPoint struct {
	ID          int    `db:"id"`
	Name        string `db:"name" json:"name"`
	Address     string `db:"address" json:"address"`
	PhoneNumber string `db:"phone_number" json:"phone_number"`
}
