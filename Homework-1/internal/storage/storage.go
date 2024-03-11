package storage

import (
	"Homework-1/internal/model"
	"bufio"
	"encoding/json"
	"io"
	"os"
)

const storageName = "storage"

type Storage struct {
	storage *os.File
	orders  []OrderDTO
}

func loadFile(file *os.File) ([]OrderDTO, error) {
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

func New() (Storage, error) {
	file, err := os.OpenFile(storageName, os.O_CREATE, 0777)
	if err != nil {
		return Storage{}, err
	}
	ords, err := loadFile(file)
	return Storage{storage: file, orders: ords}, nil

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

func (s *Storage) SaveChanges() error {
	err := writeBytes(s.orders)
	if err != nil {
		return err
	}
	return nil
}

// Create creates order
func (s *Storage) Create(input model.OrderInput) error {
	newOrder := OrderDTO{
		ID:         input.ID,
		CustomerID: input.CustomerID,
		ExpireDate: input.ExpireDate,
	}
	s.orders = append(s.orders, newOrder)
	err := s.SaveChanges()
	if err != nil {
		return err
	}
	return nil
}

// GetAllOrders returns all orders
func (s *Storage) GetAllOrders() []OrderDTO {
	return s.orders
}
