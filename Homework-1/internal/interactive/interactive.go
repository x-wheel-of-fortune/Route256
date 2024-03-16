package interactive

import (
	"Homework-1/internal/model"
	"fmt"
	"strconv"
	"sync"
)

type Storage struct {
	points map[int]model.PickupPoint
	mx     *sync.RWMutex
}

func NewStorage() *Storage {
	points := make(map[int]model.PickupPoint)
	mx := &sync.RWMutex{}
	return &Storage{points: points, mx: mx}
}

func (s Storage) Add(id int, name string, address string, phoneNumber string) {
	newPoint := model.PickupPoint{
		ID:          id,
		Name:        name,
		Address:     address,
		PhoneNumber: phoneNumber,
	}
	s.mx.Lock()
	defer s.mx.Unlock()
	s.points[id] = newPoint
	fmt.Println("ПВЗ", newPoint, "добавлен в базу")
}

func (s Storage) PrintOne(id int) model.PickupPoint {
	s.mx.RLock()
	defer s.mx.RUnlock()
	point := s.points[id]
	fmt.Println("Информация о ПВЗ с id=" + strconv.Itoa(id) + ":")
	fmt.Println(point)
	return point
}

func (s Storage) PrintAll() {
	s.mx.RLock()
	defer s.mx.RUnlock()
	points := s.points
	fmt.Println(points)
}

func Run() {
	var id int
	var name, address, phoneNumber string

	stor := NewStorage()
	for {
		fmt.Println("Введите команду")
		fmt.Println("1 - Добавить ПВЗ")
		fmt.Println("2 - Информация о ПВЗ по id")
		fmt.Println("3 - Список всех ПВЗ")
		fmt.Println("4 - Запустить добавление и чтение одновременно")
		var command int
		_, err := fmt.Scanf("%d", &command)
		if err != nil {
			fmt.Println(err)
		}
		switch command {
		case 1:
			fmt.Println("Введите id добавляемого ПВЗ")
			_, err = fmt.Scanf("%d", &id)
			fmt.Println("Введите название добавляемого ПВЗ")
			_, err = fmt.Scanf("%s", &name)
			fmt.Println("Введите адрес добавляемого ПВЗ")
			_, err = fmt.Scanf("%s", &address)
			fmt.Println("Введите контактный номер добавляемого ПВЗ")
			_, err = fmt.Scanf("%s", &phoneNumber)

			fmt.Println(id, name, address, phoneNumber)
			go stor.Add(id, name, address, phoneNumber)

		case 2:
			fmt.Println("Введите id ПВЗ")
			_, err = fmt.Scanf("%d", &id)
			go stor.PrintOne(id)

		case 3:
			go stor.PrintAll()

		case 4:
			for i := 1; i < 10; i++ {
				go stor.Add(i, "Название_"+strconv.Itoa(i), "Адрес_"+strconv.Itoa(i), "Номер_"+strconv.Itoa(i))
				go stor.PrintAll()
			}

		default:
			fmt.Println("Некорректная команда")

		}
	}

}
