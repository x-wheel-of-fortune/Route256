package kafka

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

// MockProducer is a mock implementation of the Kafka producer interface
type MockProducer struct {
	mock.Mock
}

// SendMessage simulates sending a message asynchronously
func (m *MockProducer) SendMessage(message *sarama.ProducerMessage) {
	m.Called(message)
}

type MockConsumer struct {
	mock.Mock
}

// ConsumePartition simulates consuming messages from a partition
func (m *MockConsumer) ConsumePartition(topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {
	args := m.Called(topic, partition, offset)
	return args.Get(0).(sarama.PartitionConsumer), args.Error(1)
}

func TestSendMessage(t *testing.T) {
	t.Parallel()

	// Create a mock producer
	mockProducer := new(MockProducer)

	// Set up expectations
	mockProducer.On("SendAsyncMessage", mock.Anything).Return()

	// Create a KafkaSender with the mock producer
	sender := NewKafkaSender(mockProducer, "test_topic")

	// Prepare a test message
	message := InfoMessage{
		Timestamp: time.Now(),
		Method:    "GET",
		Raw:       []byte("test message"),
	}

	// Call the method under test
	err := sender.sendAsyncMessage(message)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Verify that the mock producer's SendAsyncMessage method was called with the correct arguments
	mockProducer.AssertCalled(t, "SendAsyncMessage", mock.Anything)
}
