package answer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"homework/internal/infrastructure/kafka"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Sender interface {
	SendMessage(message InfoMessage) error
}

type InfoMessage struct {
	Timestamp time.Time
	Method    string
	Raw       []byte
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

type UpdatePickupPointRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

func (s *KafkaSender) SendMessage(message InfoMessage) error {
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

func (s *KafkaSender) buildMessage(message InfoMessage) (*sarama.ProducerMessage, error) {
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

func AuthMiddleware(handler http.Handler, sender Sender) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		true_username, username_exists := os.LookupEnv("USER")
		true_password, password_exists := os.LookupEnv("PASSWORD")
		if username_exists && password_exists {
			username, password, ok := req.BasicAuth()
			if !ok || username != true_username || password != true_password {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		if req.Method == http.MethodDelete || req.Method == http.MethodGet {
			sender.SendMessage(InfoMessage{
				Timestamp: time.Now(),
				Method:    req.Method,
				Raw:       []byte(fmt.Sprintf("%s", req.URL)),
			})
		} else {
			body, err := io.ReadAll(req.Body)

			sender.SendMessage(InfoMessage{
				Timestamp: time.Now(),
				Method:    req.Method,
				Raw:       body,
			})

			req.Body.Close() //  must close
			req.Body = io.NopCloser(bytes.NewBuffer(body))
			var unm UpdatePickupPointRequest
			if err = json.Unmarshal(body, &unm); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		handler.ServeHTTP(w, req)

	}
}
