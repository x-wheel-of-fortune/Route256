package answer

import (
	"Homework-1/internal/infrastructure/kafka"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"time"
)

type Sender interface {
	sendAsyncMessage(message infoMessage) error
	SendMessage(message infoMessage) error
	sendMessages(messages []infoMessage) error
}

type infoMessage struct {
	Timestamp time.Time
	Method    string
	Raw       string
}

type KafkaSender struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaSender(producer *kafka.Producer, topic string) *KafkaSender {
	return &KafkaSender{
		producer,
		topic,
	}
}

func (s *KafkaSender) sendAsyncMessage(message infoMessage) error {
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		fmt.Println("Send message marshal error", err)
		return err
	}

	s.producer.SendAsyncMessage(kafkaMsg)

	fmt.Println("Send async message with key:", kafkaMsg.Key)
	return nil
}

func (s *KafkaSender) sendMessage(message infoMessage) error {
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		fmt.Println("Send message marshal error", err)
		return err
	}

	//partition, offset, err := s.producer.SendSyncMessage(kafkaMsg)
	_, _, err = s.producer.SendSyncMessage(kafkaMsg)

	if err != nil {
		fmt.Println("Send message connector error", err)
		return err
	}

	//fmt.Println("Partition: ", partition, " Offset: ", offset, " AnswerID:", message.AnswerID)
	return nil
}

func (s *KafkaSender) sendMessages(messages []infoMessage) error {
	var kafkaMsg []*sarama.ProducerMessage
	var message *sarama.ProducerMessage
	var err error

	for _, m := range messages {
		message, err = s.buildMessage(m)
		kafkaMsg = append(kafkaMsg, message)

		if err != nil {
			fmt.Println("Send message marshal error", err)
			return err
		}
	}

	err = s.producer.SendSyncMessages(kafkaMsg)

	if err != nil {
		fmt.Println("Send message connector error", err)
		return err
	}

	fmt.Println("Send messages count:", len(messages))
	return nil
}

func (s *KafkaSender) buildMessage(message infoMessage) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(message)

	if err != nil {
		fmt.Println("Send message marshal error", err)
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic:     s.topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
		Key:       sarama.StringEncoder(fmt.Sprint(message.Raw)),
		Headers: []sarama.RecordHeader{ // например, в хедер можно записать версию релиза
			{
				Key:   []byte("test-header"),
				Value: []byte("test-value"),
			},
		},
	}, nil
}
