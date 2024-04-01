package model

import (
	"time"
)

type PackagingType string

type Order struct {
	ID                 int           `json:"order_id"`
	CustomerID         int           `json:"customer_id"`
	ExpireDate         time.Time     `json:"expire_date"`
	IsFinished         bool          `json:"is_finished"`
	DateFinished       time.Time     `json:"date_finished"`
	IsReturnedByClient bool          `json:"is_returned_by_client"`
	IsDeleted          bool          `json:"is_deleted"`
	Weight             float64       `json:"weight"`
	Price              int           `json:"price"`
	Packaging          PackagingType `json:"packaging"`
}

type OrderInput struct {
	ID                 int       `json:"order_id"`
	CustomerID         int       `json:"customer_id"`
	ExpireDate         time.Time `json:"expire_date"`
	IsFinished         bool      `json:"is_finished"`
	DateFinished       time.Time `json:"date_finished"`
	IsReturnedByClient bool      `json:"is_returned_by_client"`
	IsDeleted          bool      `json:"is_deleted"`
	Weight             float64   `json:"weight"`
	Price              int       `json:"price"`
	Packaging          string    `json:"packaging"`
}
