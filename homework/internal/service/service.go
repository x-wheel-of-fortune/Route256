package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"homework/internal/model"
	"homework/internal/service/packaging"
	storage2 "homework/internal/storage"
)

type storage interface {
	Create(order model.Order) error
	GetAllOrders() []storage2.OrderDTO
	SaveChanges() error
}

type Service struct {
	storage           storage
	packagingVariants map[model.PackagingType]packaging.PackagingVariant
}

func New(s storage, pkgVar map[model.PackagingType]packaging.PackagingVariant) Service {
	return Service{
		storage:           s,
		packagingVariants: pkgVar,
	}
}

func (s *Service) processPackaging(order model.Order) (model.Order, error) {
	v := s.packagingVariants[order.Packaging]
	err := v.ValidateWeight(order.Weight)
	if err != nil {
		return model.Order{}, err
	}
	packagingExpense, err := v.CalculatePackagingExpense(order)
	if err != nil {
		return model.Order{}, err
	}
	order.Price += packagingExpense
	return order, nil
}

func (s Service) validateOrderInput(input model.OrderInput) (model.Order, error) {
	if input.ID == 0 {
		return model.Order{}, errors.New("не указан id заказа")
	}
	if input.CustomerID == 0 {
		return model.Order{}, errors.New("не указан id получателя")
	}
	if input.ExpireDateStr == "" {
		return model.Order{}, errors.New("не указан срок хранения заказа")
	}
	if input.Weight == 0 {
		return model.Order{}, errors.New("не указан вес зкаказа")
	}
	if input.Price == 0.0 {
		return model.Order{}, errors.New("не указана стоимость заказа")
	}
	if input.Packaging == "" {
		return model.Order{}, errors.New("не указана форма упаковки заказа")
	}
	_, exists := s.packagingVariants[model.PackagingType(input.Packaging)]
	if !exists {
		return model.Order{}, errors.New("некорректная форма упаковки заказа")
	}
	expireDate, err := time.Parse("2006-1-2", input.ExpireDateStr)
	if err != nil {
		return model.Order{}, err
	}
	if expireDate.Before(time.Now()) {
		return model.Order{}, errors.New("срок хранения товара находится в прошлом")
	}

	orders := s.storage.GetAllOrders()
	for _, order := range orders {
		if input.ID == order.ID {
			return model.Order{}, errors.New("заказ с этим id уже есть в базе")
		}
	}

	newOrder := model.Order{
		ID:         input.ID,
		CustomerID: input.CustomerID,
		ExpireDate: expireDate,
		Weight:     input.Weight,
		Price:      input.Price,
		Packaging:  model.PackagingType(input.Packaging),
	}

	return newOrder, nil
}

// Create ...
func (s Service) Create(input model.OrderInput) error {
	newOrder, err := s.validateOrderInput(input)
	if err != nil {
		return err
	}
	newOrder, err = s.processPackaging(newOrder)
	if err != nil {
		return err
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
				Weight:             order.Weight,
				Price:              order.Price,
				Packaging:          model.PackagingType(order.Packaging),
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
				Weight:             order.Weight,
				Price:              order.Price,
				Packaging:          model.PackagingType(order.Packaging),
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
