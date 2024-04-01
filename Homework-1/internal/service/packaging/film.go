package packaging

import (
	"Homework-1/internal/model"
)

const (
	filmMaxWeight float64 = 1.7e+308
	filmPrice     int     = 1
)

type FilmPackaging struct {
}

func (v FilmPackaging) GetMaxWeight() float64 {
	return filmMaxWeight
}

func (v FilmPackaging) GetPrice() int {
	return filmPrice
}

func (v FilmPackaging) validateWeight(weight float64) error {
	return nil
}

func (v FilmPackaging) calculatePackagingExpense(order model.Order) (int, error) {
	// Пока что функция никак не использует полученный на вход order, но в будущем
	// логика вычисления стоимости упаковки может учитывать значения некоторых полей
	// обрабатываемого заказа
	expense := v.GetPrice()
	return expense, nil
}

func (v FilmPackaging) ProcessPackaging(order model.Order) (model.Order, error) {
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
