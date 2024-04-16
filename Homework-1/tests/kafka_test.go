//go:build integration
// +build integration

package tests

import (
	"Homework-1/internal/infrastructure/answer"
	"Homework-1/internal/infrastructure/info"
	"Homework-1/internal/infrastructure/kafka"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func cleanupTopic(topic string, broker string) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0

	// Create a new Kafka admin client
	admin, err := sarama.NewClusterAdmin([]string{broker}, config)
	if err != nil {
		return err
	}
	defer admin.Close()

	// Delete the topic
	err = admin.DeleteTopic(topic)
	if err != nil {
		return err
	}

	return nil
}

func TestPOST(t *testing.T) {
	var (
		timestamp = time.Unix(123, 123)
		body      = []byte("{\n    \"name\": \"PickupPoint_1\",\n    \"address\":\"Address_1\",\n    \"phone_number\":\"+7-999-999-99-99\"\n  }")
	)
	// Create a unique topic for this test
	topic := "test_POST_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	t.Run("smoke test", func(t *testing.T) {

		// arrange
		broker, exists := os.LookupEnv("TEST_KAFKA_BROKER")
		require.True(t, exists)
		kafkaProducer, err := kafka.NewProducer([]string{broker})
		defer kafkaProducer.Close()
		sender := answer.NewKafkaSender(kafkaProducer, topic)

		kafkaConsumer, err := kafka.NewConsumer([]string{broker})
		handlers := map[string]info.HandleFunc{
			topic: func(message *sarama.ConsumerMessage) {
				pm := answer.InfoMessage{}
				err = json.Unmarshal(message.Value, &pm)
				require.NoError(t, err)
			},
		}

		infos := info.NewService(info.NewReceiver(kafkaConsumer, handlers))
		infoChan := make(chan string)

		// act
		infos.StartConsume(topic, infoChan)
		sender.SendMessage(answer.InfoMessage{
			Timestamp: timestamp,
			Method:    http.MethodPost,
			Raw:       body,
		})
		result := <-infoChan
		cleanupTopic(broker, topic)

		//assert
		assert.Equal(t, "Datetime:1970-01-01 07:02:03.000000123 +0700 +07\nMethod:POST\nRaw:{\n    \"name\": \"PickupPoint_1\",\n    \"address\":\"Address_1\",\n    \"phone_number\":\"+7-999-999-99-99\"\n  }\n\n", result)
	})

}

func TestPUT(t *testing.T) {
	var (
		timestamp = time.Unix(123, 123)
		body      = []byte("{\n    \"id\":1,\n    \"name\": \"Updated_PickupPoint_1\",\n    \"address\":\"Updated_Address_1\",\n    \"phone_number\":\"+7-999-999-99-99\"\n}")
	)
	topic := "test_PUT_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	t.Run("smoke test", func(t *testing.T) {

		// arrange
		broker, exists := os.LookupEnv("TEST_KAFKA_BROKER")
		require.True(t, exists)
		kafkaProducer, err := kafka.NewProducer([]string{broker})
		defer kafkaProducer.Close()
		sender := answer.NewKafkaSender(kafkaProducer, topic)

		kafkaConsumer, err := kafka.NewConsumer([]string{broker})
		handlers := map[string]info.HandleFunc{
			topic: func(message *sarama.ConsumerMessage) {
				pm := answer.InfoMessage{}
				err = json.Unmarshal(message.Value, &pm)
				require.NoError(t, err)
			},
		}

		infos := info.NewService(info.NewReceiver(kafkaConsumer, handlers))
		infoChan := make(chan string)

		// act
		infos.StartConsume(topic, infoChan)
		sender.SendMessage(answer.InfoMessage{
			Timestamp: timestamp,
			Method:    http.MethodPut,
			Raw:       body,
		})
		result := <-infoChan
		cleanupTopic(broker, topic)

		//assert
		assert.Equal(t, "Datetime:1970-01-01 07:02:03.000000123 +0700 +07\nMethod:PUT\nRaw:{\n    \"id\":1,\n    \"name\": \"Updated_PickupPoint_1\",\n    \"address\":\"Updated_Address_1\",\n    \"phone_number\":\"+7-999-999-99-99\"\n}\n\n", result)
	})

}

func TestGET(t *testing.T) {
	var (
		timestamp = time.Unix(123, 123)
		url       = "/pickup_point/1"
	)
	topic := "test_GET_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	t.Run("smoke test", func(t *testing.T) {

		// arrange
		broker, exists := os.LookupEnv("TEST_KAFKA_BROKER")
		require.True(t, exists)
		kafkaProducer, err := kafka.NewProducer([]string{broker})
		defer kafkaProducer.Close()
		sender := answer.NewKafkaSender(kafkaProducer, topic)

		kafkaConsumer, err := kafka.NewConsumer([]string{broker})
		handlers := map[string]info.HandleFunc{
			topic: func(message *sarama.ConsumerMessage) {
				pm := answer.InfoMessage{}
				err = json.Unmarshal(message.Value, &pm)
				require.NoError(t, err)
			},
		}

		infos := info.NewService(info.NewReceiver(kafkaConsumer, handlers))
		infoChan := make(chan string)

		// act
		infos.StartConsume(topic, infoChan)
		sender.SendMessage(answer.InfoMessage{
			Timestamp: timestamp,
			Method:    http.MethodGet,
			Raw:       []byte(url),
		})
		result := <-infoChan

		//assert
		assert.Equal(t, "Datetime:1970-01-01 07:02:03.000000123 +0700 +07\nMethod:GET\nRaw:/pickup_point/1\n\n", result)
	})

}

func TestDELETE(t *testing.T) {
	var (
		timestamp = time.Unix(123, 123)
		url       = "/pickup_point/1"
	)
	topic := "test_DELETE_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	t.Run("smoke test", func(t *testing.T) {

		// arrange
		broker, exists := os.LookupEnv("TEST_KAFKA_BROKER")
		require.True(t, exists)
		kafkaProducer, err := kafka.NewProducer([]string{broker})
		defer kafkaProducer.Close()
		sender := answer.NewKafkaSender(kafkaProducer, topic)

		kafkaConsumer, err := kafka.NewConsumer([]string{broker})
		handlers := map[string]info.HandleFunc{
			topic: func(message *sarama.ConsumerMessage) {
				pm := answer.InfoMessage{}
				err = json.Unmarshal(message.Value, &pm)
				require.NoError(t, err)
			},
		}

		infos := info.NewService(info.NewReceiver(kafkaConsumer, handlers))
		infoChan := make(chan string)

		// act
		infos.StartConsume(topic, infoChan)
		sender.SendMessage(answer.InfoMessage{
			Timestamp: timestamp,
			Method:    http.MethodDelete,
			Raw:       []byte(url),
		})
		result := <-infoChan

		//assert
		assert.Equal(t, "Datetime:1970-01-01 07:02:03.000000123 +0700 +07\nMethod:DELETE\nRaw:/pickup_point/1\n\n", result)
	})

}
