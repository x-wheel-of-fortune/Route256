package info

import (
	"Homework-1/internal/infrastructure/answer"
	"Homework-1/internal/infrastructure/kafka"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
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

	// получаем все партиции топика
	partitionList, err := r.consumer.SingleConsumer.Partitions(topic)

	if err != nil {
		return err
	}

	/*
	   sarama.OffsetOldest - перечитываем каждый раз все
	   sarama.OffsetNewest - перечитываем только новые

	   Можем задавать отдельно на каждую партицию
	   Также можем сходить в отдельное хранилище и взять оттуда сохраненный offset
	*/
	initialOffset := sarama.OffsetOldest

	for _, partition := range partitionList {
		pc, err := r.consumer.SingleConsumer.ConsumePartition(topic, partition, initialOffset)

		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer, partition int32) {
			for message := range pc.Messages() {
				handler(message)
				//fmt.Println("Read Topic: ", topic, " Partition: ", partition, " Offset: ", message.Offset)
				//fmt.Println("Received Key: ", string(message.Key), " Value: ", string(message.Value))
				var unm answer.InfoMessage
				json.Unmarshal(message.Value, &unm)
				//fmt.Printf("Datetime:%v\nMethod:%s\nRaw:%s\n\n", unm.Timestamp, unm.Method, string(unm.Raw))
				infoChan <- fmt.Sprintf("Datetime:%v\nMethod:%s\nRaw:%s\n\n", unm.Timestamp, unm.Method, string(unm.Raw))

			}
		}(pc, partition)
	}

	return nil
}
