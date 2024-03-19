package interactive_storage

import (
	"Homework-1/internal/model"
	"errors"
	"fmt"
	"sync"
)

type Storage struct {
	points map[int]model.PickupPoint
	mx     sync.RWMutex
}

func NewStorage() *Storage {
	points := make(map[int]model.PickupPoint)
	mx := sync.RWMutex{}
	return &Storage{points: points, mx: mx}
}

func (s Storage) Add(newPoint model.PickupPoint) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	_, ok := s.points[newPoint.ID]
	if ok {
		return errors.New(fmt.Sprintf("ПВЗ с id=%d уже существует", newPoint.ID))
	}
	s.points[newPoint.ID] = newPoint
	return nil
}

func (s Storage) Get(id int) (model.PickupPoint, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	point, ok := s.points[id]
	if !ok {
		return model.PickupPoint{}, errors.New(fmt.Sprintf("ПВЗ с id=%d не найден", id))
	}
	return point, nil
}
