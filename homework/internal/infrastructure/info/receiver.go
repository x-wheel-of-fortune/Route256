package info

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"

	"homework/internal/infrastructure/answer"
	"homework/internal/infrastructure/kafka"
)

type HandleFunc func(message *sarama.ConsumerMessage)

type KafkaReceiver struct {
	consumer *kafka.Consumer
	handlers map[string]HandleFunc
}

func NewReceiver(consumer *kafka.Consumer, handlers map[string]HandleFunc) *KafkaReceiver {
	return &KafkaReceiver{
		consumer: consumer,
		handlers: handlers,
	}
}

func (r *KafkaReceiver) Subscribe(topic string, infoChan chan<- string) error {
	handler, ok := r.handlers[topic]

	if !ok {
		return errors.New("can not find handler")
	}

	partitionList, err := r.consumer.SingleConsumer.Partitions(topic)

	if err != nil {
		return err
	}

	initialOffset := sarama.OffsetOldest

	for _, partition := range partitionList {
		pc, err := r.consumer.SingleConsumer.ConsumePartition(topic, partition, initialOffset)

		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer, partition int32) {
			for message := range pc.Messages() {
				handler(message)
				var unm answer.InfoMessage
				json.Unmarshal(message.Value, &unm)
				infoChan <- fmt.Sprintf("Datetime:%v\nMethod:%s\nRaw:%s\n\n", unm.Timestamp, unm.Method, string(unm.Raw))

			}
		}(pc, partition)
	}

	return nil
}
