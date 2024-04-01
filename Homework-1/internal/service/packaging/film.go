package packaging

import (
	"Homework-1/internal/model"
)

type FilmPackaging struct {
}

func (v FilmPackaging) ProcessPackaging(order model.Order) (model.Order, error) {
	order.Price += 1
	return order, nil
}
