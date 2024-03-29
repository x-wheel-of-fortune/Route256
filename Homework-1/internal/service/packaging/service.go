package packaging

import "Homework-1/internal/model"

const (
	BagPackagingType  string = "bag"
	BoxPackagingType  string = "box"
	FilmPackagingType string = "film"
)

type PackagingVariant interface {
	ProcessPackaging(order model.OrderInput) (model.OrderInput, error)
}
