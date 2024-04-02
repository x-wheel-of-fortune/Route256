package packaging

import "Homework-1/internal/model"

const (
	BagPackagingType  model.PackagingType = "bag"
	BoxPackagingType  model.PackagingType = "box"
	FilmPackagingType model.PackagingType = "film"
)

type PackagingVariant interface {
	GetMaxWeight() float64
	GetPrice() int
	ValidateWeight(weight float64) error
	CalculatePackagingExpense(order model.Order) (int, error)
}
