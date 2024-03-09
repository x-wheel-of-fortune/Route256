package service

import (
	"Homework-1/internal/model"
	"errors"
	"fmt"
)

type storage interface {
	Create(order model.OrderInput) error
	Delete(id int) error
	List(customerId int, limit int, onlyNotFinished bool) ([]model.Order, error)
	Returns() ([]model.Order, error)
	Finish(ids []int) error
	Return(id int, customerId int) error
}

type Service struct {
	s storage
}

func New(s storage) Service {
	return Service{s: s}
}

// Create ...
func (s Service) Create(input model.OrderInput) error {
	//if len(input.Description) == 0 {
	//	return errors.New("пустое описание цели")
	//}

	return s.s.Create(input)
}

// Delete ...
func (s Service) Delete(id int) error {
	if id == 0 {
		return errors.New("нулевой id цели")
	}
	return s.s.Delete(id)
}

// List ...
func (s Service) List(customerId int, limit int, onlyNotFinished bool) ([]model.Order, error) {
	orders, err := s.s.List(customerId, limit, onlyNotFinished)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// Return ...
func (s Service) Return(id int, customerId int) error {
	return s.s.Return(id, customerId)
}

// Returns ...
func (s Service) Returns(resultsPerPage int) (string, error) {
	orders, err := s.s.Returns()
	if err != nil {
		return "", err
	}

	paginatedReturns := ""
	count := 0
	for _, order := range orders {
		if count%resultsPerPage == 0 {
			paginatedReturns += fmt.Sprintf("Страница %d\n", (count/resultsPerPage + 1))
		}
		s := fmt.Sprintf("id заказа: %d, ", order.ID)
		s += fmt.Sprintf("выдан, дата выдачи: %d-%d-%d", order.DateFinished.Year(), order.DateFinished.Month(), order.DateFinished.Day())
		s += ", клиент оформил возврат"

		fmt.Println(s)
		paginatedReturns += s + "\n"

		count += 1
	}
	return paginatedReturns, nil
}

func (s Service) Finish(ids []int) error {
	err := s.s.Finish(ids)
	if err != nil {
		return err
	}
	return nil
}
