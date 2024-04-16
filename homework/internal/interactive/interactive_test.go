package interactive

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"homework/internal/model"
	"testing"
)

func Test_Adder(t *testing.T) {
	t.Parallel()
	var (
		point = model.PickupPoint{
			ID:          1,
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		addChan := make(chan model.PickupPoint)
		infoChan := make(chan string)
		s := setUp(t)
		defer s.tearDown()
		s.mockInteractiveStorage.EXPECT().Add(point).Return(nil)
		go Adder(addChan, s.mockInteractiveStorage, infoChan)
		addChan <- point
		result := <-infoChan
		assert.Equal(t, "Adder начал добавление ПВЗ с id=1 в базу", result)
		result = <-infoChan
		assert.Equal(t, "ПВЗ с id=1 успешно добавлен в базу", result)
		result = <-infoChan
		assert.Equal(t, "Adder закончил добавление ПВЗ с id=1 в базу", result)
	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()
		addChan := make(chan model.PickupPoint)
		infoChan := make(chan string)
		s := setUp(t)
		defer s.tearDown()
		s.mockInteractiveStorage.EXPECT().Add(point).Return(errors.New("internal error"))
		go Adder(addChan, s.mockInteractiveStorage, infoChan)
		addChan <- point
		result := <-infoChan
		assert.Equal(t, "Adder начал добавление ПВЗ с id=1 в базу", result)
		result = <-infoChan
		assert.Equal(t, "internal error", result)
		result = <-infoChan
		assert.Equal(t, "Adder закончил добавление ПВЗ с id=1 в базу", result)
	})

}

func Test_Reader(t *testing.T) {
	t.Parallel()
	var (
		point = model.PickupPoint{
			ID:          1,
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		id = 1
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		readChan := make(chan int)
		infoChan := make(chan string)
		s := setUp(t)
		defer s.tearDown()
		s.mockInteractiveStorage.EXPECT().Get(id).Return(point, nil)
		go Reader(readChan, s.mockInteractiveStorage, infoChan)
		readChan <- id
		result := <-infoChan
		assert.Equal(t, "Reader начал чтение ПВЗ с id=1", result)
		result = <-infoChan
		assert.Equal(t, "{ID:1 Name:Name Address:Address PhoneNumber:PhoneNumber}\n", result)
		result = <-infoChan
		assert.Equal(t, "Reader закончил чтение ПВЗ с id=1", result)
	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()
		readChan := make(chan int)
		infoChan := make(chan string)
		s := setUp(t)
		defer s.tearDown()
		s.mockInteractiveStorage.EXPECT().Get(id).Return(model.PickupPoint{}, errors.New("internal error"))
		go Reader(readChan, s.mockInteractiveStorage, infoChan)
		readChan <- id
		result := <-infoChan
		assert.Equal(t, "Reader начал чтение ПВЗ с id=1", result)
		result = <-infoChan
		assert.Equal(t, "internal error", result)
		result = <-infoChan
		assert.Equal(t, "Reader закончил чтение ПВЗ с id=1", result)
	})

}
