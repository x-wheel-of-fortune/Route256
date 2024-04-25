package packaging

import (
	"homework/internal/model"
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

func (v FilmPackaging) ValidateWeight(weight float64) error {
	return nil
}

func (v FilmPackaging) CalculatePackagingExpense(order model.Order) (int, error) {
	// Пока что функция никак не использует полученный на вход order, но в будущем
	// логика вычисления стоимости упаковки может учитывать значения некоторых полей
	// обрабатываемого заказа
	expense := v.GetPrice()
	return expense, nil
}
