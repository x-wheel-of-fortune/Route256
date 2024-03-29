package packaging

import (
	"Homework-1/internal/model"
	"errors"
)

type BagPackaging struct {
}

func (v BagPackaging) ProcessPackaging(order model.OrderInput) (model.OrderInput, error) {
	if order.Weight >= 10 {
		return model.OrderInput{}, errors.New("в пакет можно упаковывать только заказы весом до 10кг")
	}
	order.Price += 5
	return order, nil
}
