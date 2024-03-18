package interactive

import (
	"Homework-1/internal/model"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
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

func (s Storage) PrintAll() {
	s.mx.RLock()
	defer s.mx.RUnlock()
	points := s.points
	fmt.Println(points)
}

func Adder(addChanel chan model.PickupPoint, stor *Storage, infoChannel chan string) {
	for {
		newPoint, ok := <-addChanel
		if !ok {
			return
		}
		infoChannel <- fmt.Sprintf("Adder начал добавление ПВЗ с id=%d в базу", newPoint.ID)
		time.Sleep(5 * time.Second)
		err := stor.Add(newPoint)
		if err != nil {
			//infoChannel <- err.Error()
			fmt.Println(err.Error())
		} else {
			fmt.Println(fmt.Sprintf("ПВЗ с id=%d успешно добавлен в базу", newPoint.ID))
		}
		infoChannel <- fmt.Sprintf("Adder закончил добавление ПВЗ с id=%d в базу", newPoint.ID)
	}
}

func Reader(readChanel chan int, stor *Storage, infoChannel chan string) {
	for {
		id, ok := <-readChanel
		if !ok {
			return
		}
		infoChannel <- fmt.Sprintf("Reader начал чтение ПВЗ с id=%d", id)
		time.Sleep(2 * time.Second)
		point, err := stor.Get(id)
		if err != nil {
			//infoChannel <- err.Error()
			fmt.Println(err.Error())
		} else {
			fmt.Printf("Информация о ПВЗ c id=%d: ", id)
			fmt.Printf("%+v\n", point)
		}
		infoChannel <- fmt.Sprintf("Reader закончил чтение ПВЗ с id=%d", id)
	}
}

func Overseer(infoChannel chan string) {
	for {
		info, ok := <-infoChannel
		if !ok {
			return
		}
		log.Println(info)
	}
}

func Run(ctx context.Context) {
	var id int
	var name, address, phoneNumber string
	stor := NewStorage()

	infoChannel := make(chan string)
	defer close(infoChannel)
	go Overseer(infoChannel)

	addChannel := make(chan model.PickupPoint)
	defer close(addChannel)
	go Adder(addChannel, stor, infoChannel)

	readChannel := make(chan int)
	defer close(readChannel)
	go Reader(readChannel, stor, infoChannel)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("\nВведите команду")
			fmt.Println("1 - Добавить ПВЗ")
			fmt.Println("2 - Узнать информацию о ПВЗ")
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

				newPoint := model.PickupPoint{ID: id, Name: name, Address: address, PhoneNumber: phoneNumber}
				addChannel <- newPoint

			case 2:
				fmt.Println("Введите id ПВЗ, информацию о котором вы хотите узнать")
				_, err = fmt.Scanf("%d", &id)
				readChannel <- id

			default:
				fmt.Println("Некорректная команда")
			}

		}
	}

}
