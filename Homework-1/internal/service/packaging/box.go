package packaging

import (
	"Homework-1/internal/model"
	"errors"
	"fmt"
)

const (
	boxMaxWeight float64 = 30
	boxPrice     int     = 20
)

type BoxPackaging struct {
}

func (v BoxPackaging) GetMaxWeight() float64 {
	return boxMaxWeight
}

func (v BoxPackaging) GetPrice() int {
	return boxPrice
}

func (v BoxPackaging) validateWeight(weight float64) error {
	if weight >= v.GetMaxWeight() {
		return errors.New(fmt.Sprintf("в коробку можно упаковывать только заказы весом до %dкг", v.GetMaxWeight()))
	}
	return nil
}

func (v BoxPackaging) calculatePackagingExpense(order model.Order) (int, error) {
	// Пока что функция никак не использует полученный на вход order, но в будущем
	// логика вычисления стоимости упаковки может учитывать значения некоторых полей
	// обрабатываемого заказа
	expense := v.GetPrice()
	return expense, nil
}

func (v BoxPackaging) ProcessPackaging(order model.Order) (model.Order, error) {
	err := v.validateWeight(order.Weight)
	if err != nil {
		return model.Order{}, err
	}
	packagingExpense, err := v.calculatePackagingExpense(order)
	if err != nil {
		return model.Order{}, err
	}
	order.Price += packagingExpense
	return order, nil
}
