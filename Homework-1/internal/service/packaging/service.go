package packaging

import "Homework-1/internal/model"

const (
	BagPackagingType  model.PackagingType = "bag"
	BoxPackagingType  model.PackagingType = "box"
	FilmPackagingType model.PackagingType = "film"
)

type PackagingVariant interface {
	ProcessPackaging(order model.OrderInput) (model.OrderInput, error)
}
