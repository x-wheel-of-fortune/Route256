package packaging

import (
	"Homework-1/internal/model"
	"errors"
)

type BagPackaging struct {
}

func (v BagPackaging) ProcessPackaging(order model.Order) (model.Order, error) {
	if order.Weight >= 10 {
		return model.Order{}, errors.New("в пакет можно упаковывать только заказы весом до 10кг")
	}
	order.Price += 5
	return order, nil
}
