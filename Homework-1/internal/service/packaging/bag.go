package packaging

import (
	"Homework-1/internal/model"
	"errors"
	"fmt"
)

const (
	bagMaxWeight float64 = 10
	bagPrice     int     = 5
)

type BagPackaging struct {
}

func (v BagPackaging) GetMaxWeight() float64 {
	return bagMaxWeight
}

func (v BagPackaging) GetPrice() int {
	return bagPrice
}

func (v BagPackaging) ValidateWeight(weight float64) error {
	if weight >= v.GetMaxWeight() {
		return errors.New(fmt.Sprintf("в пакет можно упаковывать только заказы весом до %dкг", int(v.GetMaxWeight())))
	}
	return nil
}

func (v BagPackaging) CalculatePackagingExpense(order model.Order) (int, error) {
	// Пока что функция никак не использует полученный на вход order, но в будущем
	// логика вычисления стоимости упаковки может учитывать значения некоторых полей
	// обрабатываемого заказа
	expense := v.GetPrice()
	return expense, nil
}
