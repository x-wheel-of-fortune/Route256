package storage

import (
	"Homework-1/internal/model"
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"
)

const storageName = "storage"

type Storage struct {
	storage *os.File
	orders  []OrderDTO
}

func New() (Storage, error) {
	file, err := os.OpenFile(storageName, os.O_CREATE, 0777)
	if err != nil {
		return Storage{}, err
	}
	ords, err := listAll(file)
	return Storage{storage: file, orders: ords}, nil

}

// Create creates order
func (s *Storage) Create(input model.OrderInput) error {

	for _, i := range s.orders {
		if input.ID == i.ID {
			return errors.New("заказ с этим id уже есть в базе")
		}
	}
	newOrder := OrderDTO{
		ID:         input.ID,
		CustomerID: input.CustomerID,
		ExpireDate: input.ExpireDate,
	}

	s.orders = append(s.orders, newOrder)
	err := writeBytes(s.orders)
	if err != nil {
		return err
	}
	return nil
}

func writeBytes(orders []OrderDTO) error {
	rawBytes, err := json.Marshal(orders)
	if err != nil {
		return err
	}

	err = os.WriteFile(storageName, rawBytes, 0777)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes an order
func (s *Storage) Delete(id int) error {

	found := false
	for indx, order := range s.orders {
		if order.ID == id {
			if order.IsFinished {
				return errors.New("этот заказ уже был выдан получателю")
			}
			if order.ExpireDate.After(time.Now()) {
				return errors.New("срок хранения этого заказа ещё не истёк")
			}
			s.orders[indx].IsDeleted = true
			found = true
		}
	}
	if !found {
		return errors.New("заказ с данным id не найден")
	}
	err := writeBytes(s.orders)
	if err != nil {
		return err
	}
	return nil
}

// Return set the target order IsReturnedByClient status to True
func (s *Storage) Return(id int, customer_id int) error {

	found := false
	for indx, order := range s.orders {
		if order.ID == id {
			if order.IsReturnedByClient {
				return errors.New("возврат этого заказа уже был принят")
			}
			if order.CustomerID != customer_id {
				return errors.New("id клиента, возвращающего заказ, не совпадает с id получателя")
			}
			if !order.IsFinished {
				return errors.New("этот заказ ещё не был выдан получателю")
			}
			if time.Now().After(order.DateFinished.Add(24 * time.Hour * 2)) {
				return errors.New("прошло уже более 2-х дней с момента выдачи заказа")
			}
			s.orders[indx].IsReturnedByClient = true
			found = true
		}
	}
	if !found {
		return errors.New("заказ с данным id не найден")
	}
	err := writeBytes(s.orders)
	if err != nil {
		return err
	}
	return nil
}

// Finish finishes an order
func (s *Storage) Finish(ids []int) error {

	customerId := 0
	// TODO Использовать вложенные циклы неэффективно, переделать с использованием множества
	for _, id := range ids {
		found := false
		for _, order := range s.orders {
			if order.ID == id {
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
			}
		}
		if !found {
			return errors.New("некоторые из заказов с заданными id не найдены")
		}

	}

	for _, id := range ids {
		for indx, order := range s.orders {
			if order.ID == id {
				s.orders[indx].IsFinished = true
				s.orders[indx].DateFinished = time.Now()
			}
		}
	}

	err := writeBytes(s.orders)
	if err != nil {
		return err
	}
	return nil
}

// List returns s.orders orders of the target user from storage
func (s *Storage) List(customerId int, limit int, onlyNotFinished bool) ([]model.Order, error) {

	customer_orders := make([]model.Order, 0, len(s.orders))
	for i := len(s.orders) - 1; i >= 0; i-- {
		order := s.orders[i]
		if !order.IsDeleted && order.CustomerID == customerId && (!onlyNotFinished || !order.IsFinished) {
			customer_orders = append(customer_orders, model.Order{
				ID:                 order.ID,
				CustomerID:         order.CustomerID,
				ExpireDate:         order.ExpireDate,
				IsFinished:         order.IsFinished,
				DateFinished:       order.DateFinished,
				IsReturnedByClient: order.IsReturnedByClient,
				IsDeleted:          order.IsDeleted,
				//Description: order.Description,
			})

			if limit > 0 && len(customer_orders) == limit {
				break
			}
		}

	}

	return customer_orders, nil
}

// Returns returns s.orders returned orders from storage
func (s *Storage) Returns() ([]model.Order, error) {

	returned := make([]model.Order, 0, len(s.orders))
	for _, order := range s.orders {
		if order.IsReturnedByClient {
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

	return returned, nil
}

func listAll(file *os.File) ([]OrderDTO, error) {
	reader := bufio.NewReader(file)
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var orders []OrderDTO
	if len(rawBytes) == 0 {
		return orders, nil
	}
	err = json.Unmarshal(rawBytes, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
