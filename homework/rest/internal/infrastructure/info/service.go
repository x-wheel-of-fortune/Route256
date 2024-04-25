package info

import (
	"fmt"
)

type Receiver interface {
	Subscribe(topic string, infoChan chan<- string) error
}

type Service struct {
	receiver Receiver
}

func NewService(receiver Receiver) *Service {
	return &Service{
		receiver: receiver,
	}
}

func (s *Service) StartConsume(topic string, infoChan chan<- string) {
	err := s.receiver.Subscribe(topic, infoChan)

	if err != nil {
		infoChan <- fmt.Sprintf("Subscribe error: %v", err)
	}
}
