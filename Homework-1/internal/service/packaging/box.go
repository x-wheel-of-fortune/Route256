package packaging

import (
	"Homework-1/internal/model"
	"errors"
)

type BoxPackaging struct {
}

func (v BoxPackaging) ProcessPackaging(order model.OrderInput) (model.OrderInput, error) {
	if order.Weight >= 30 {
		return model.OrderInput{}, errors.New("в коробку можно упаковывать только заказы весом до 30кг")
	}
	order.Price += 20
	return order, nil
}
