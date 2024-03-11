package service

import (
	"Homework-1/internal/model"
	storage2 "Homework-1/internal/storage"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type storage interface {
	Create(order model.OrderInput) error
	GetAllOrders() []storage2.OrderDTO
	SaveChanges() error
}

type Service struct {
	storage storage
}

func New(s storage) Service {
	return Service{storage: s}
}

// Create ...
func (s Service) Create(orderID int, customerID int, expireDateStr string) error {
	if orderID == 0 {
		return errors.New("не указан id заказа")
	}
	if customerID == 0 {
		return errors.New("не указан id получателя")
	}
	if expireDateStr == "" {
		return errors.New("не указан срок хранения заказа")
	}
	expireDate, err := time.Parse("2006-1-2", expireDateStr)
	if err != nil {
		return err
	}
	if expireDate.Before(time.Now()) {
		return errors.New("срок хранения товара находится в прошлом")
	}

	orders := s.storage.GetAllOrders()
	for _, order := range orders {
		if orderID == order.ID {
			return errors.New("заказ с этим id уже есть в базе")
		}
	}

	newOrder := model.OrderInput{
		ID:         orderID,
		CustomerID: customerID,
		ExpireDate: expireDate,
	}

	return s.storage.Create(newOrder)
}

// Delete ...
func (s Service) Delete(id int) error {
	if id == 0 {
		return errors.New("не указано id возвращаемого заказа")
	}

	orders := s.storage.GetAllOrders()
	found := false
	for indx, order := range orders {
		if !order.IsDeleted && order.ID == id {
			if order.IsFinished {
				return errors.New("этот заказ уже был выдан получателю")
			}
			if order.ExpireDate.After(time.Now()) {
				return errors.New("срок хранения этого заказа ещё не истёк")
			}
			orders[indx].IsDeleted = true
			found = true
			break
		}
	}
	if !found {
		return errors.New("заказ с данным id не найден")
	}

	return s.storage.SaveChanges()
}

// Return ...
func (s Service) Return(id int, customerId int) error {
	if id == 0 {
		return errors.New("не указано id возвращаемого заказа")
	}
	if customerId == 0 {
		return errors.New("не указано id клиента, возвращающего заказ")
	}

	orders := s.storage.GetAllOrders()
	found := false
	for indx, order := range orders {
		if !order.IsDeleted && order.ID == id {
			if order.IsReturnedByClient {
				return errors.New("возврат этого заказа уже был принят")
			}
			if order.CustomerID != customerId {
				return errors.New("id клиента, возвращающего заказ, не совпадает с id получателя")
			}
			if !order.IsFinished {
				return errors.New("этот заказ ещё не был выдан получателю")
			}
			if time.Now().After(order.DateFinished.Add(24 * time.Hour * 2)) {
				return errors.New("прошло уже более 2-х дней с момента выдачи заказа")
			}
			orders[indx].IsReturnedByClient = true
			found = true
			break
		}
	}
	if !found {
		return errors.New("заказ с данным id не найден")
	}
	return s.storage.SaveChanges()
}

// Finish ...
func (s Service) Finish(idsStr string) error {
	if idsStr == "" {
		return errors.New("не указаны id выдаваемых заказов")
	}
	var ids []int
	err := json.Unmarshal([]byte(idsStr), &ids)
	if err != nil {
		return err
	}

	customerId := 0
	orders := s.storage.GetAllOrders()
	// TODO Использовать вложенные циклы неэффективно, переделать с использованием множества
	for _, id := range ids {
		found := false
		for _, order := range orders {
			if !order.IsDeleted && order.ID == id {
				if customerId == 0 {
					customerId = order.CustomerID
				} else {
					if customerId != order.CustomerID {
						return errors.New("не все заказы принадлежат одному клиенту")
					}
				}

				if order.IsFinished {
					return errors.New("некоторые из заказов уже были выданы клиенту")
				}

				if order.ExpireDate.Before(time.Now()) {
					return errors.New("у некоторых из заказов истёк срок хранения")
				}
				found = true
				break
			}
		}
		if !found {
			return errors.New("некоторые из заказов с заданными id не найдены")
		}

	}

	// TODO Использовать вложенные циклы неэффективно, переделать с использованием множества
	for _, id := range ids {
		for indx, order := range orders {
			if order.ID == id {
				orders[indx].IsFinished = true
				orders[indx].DateFinished = time.Now()
				break
			}
		}
	}

	return s.storage.SaveChanges()
}

// List ...
func (s Service) List(customerId int, limit int, onlyNotFinished bool) ([]model.Order, error) {
	if customerId == 0 {
		return nil, errors.New("не задано id пользователя")
	}

	orders := s.storage.GetAllOrders()
	customer_orders := make([]model.Order, 0, len(orders))
	for i := len(orders) - 1; i >= 0; i-- {
		order := orders[i]
		if !order.IsDeleted && order.CustomerID == customerId && (!onlyNotFinished || !order.IsFinished) {
			customer_orders = append(customer_orders, model.Order{
				ID:                 order.ID,
				CustomerID:         order.CustomerID,
				ExpireDate:         order.ExpireDate,
				IsFinished:         order.IsFinished,
				DateFinished:       order.DateFinished,
				IsReturnedByClient: order.IsReturnedByClient,
				IsDeleted:          order.IsDeleted,
			})

			if limit > 0 && len(customer_orders) == limit {
				break
			}
		}

	}
	return customer_orders, nil
}

// Returns ...
func (s Service) Returns(resultsPerPage int) (string, error) {
	orders := s.storage.GetAllOrders()
	returned := make([]model.Order, 0, len(orders))
	for _, order := range orders {
		if !order.IsDeleted && order.IsReturnedByClient {
			returned = append(returned, model.Order{
				ID:                 order.ID,
				CustomerID:         order.CustomerID,
				ExpireDate:         order.ExpireDate,
				IsFinished:         order.IsFinished,
				DateFinished:       order.DateFinished,
				IsReturnedByClient: order.IsReturnedByClient,
				IsDeleted:          order.IsDeleted,
			})
		}
	}

	paginatedReturns := ""
	count := 0
	for _, order := range returned {
		if count%resultsPerPage == 0 {
			paginatedReturns += fmt.Sprintf("Страница %d\n", (count/resultsPerPage + 1))
		}
		s := fmt.Sprintf("id заказа: %d, ", order.ID)
		s += fmt.Sprintf("выдан, дата выдачи: %d-%d-%d", order.DateFinished.Year(), order.DateFinished.Month(), order.DateFinished.Day())
		s += ", клиент оформил возврат"
		paginatedReturns += s + "\n"

		count += 1
	}

	return paginatedReturns, nil
}
