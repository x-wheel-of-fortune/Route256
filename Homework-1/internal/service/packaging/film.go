package packaging

import (
	"Homework-1/internal/model"
)

type FilmPackaging struct {
}

func (v FilmPackaging) ProcessPackaging(order model.OrderInput) (model.OrderInput, error) {
	order.Price += 1
	return order, nil
}
