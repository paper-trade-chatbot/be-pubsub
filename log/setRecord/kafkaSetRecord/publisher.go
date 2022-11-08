package kafkaSetRecord

import (
	"time"

	"github.com/paper-trade-chatbot/be-pubsub/kafka"
	"github.com/paper-trade-chatbot/be-pubsub/kafka/json"
)

func NewPublisher(brokers []string) (*json.PublisherImpl, error) {

	publisher, err := json.NewPublisher(&kafka.PublisherConfig{
		Config: &kafka.ConfigImpl{
			Brokers:  brokers,
			Exported: kafka.KAFKA_EXPORTED_PRIVATE,
			Domain:   []string{"log"},
			Dataset:  "set_record",
			DataType: kafka.KAFKA_DATA_TYPE_JSON,
			Version:  0,
		},
		BatchTimeout: time.Millisecond * 10,
	})

	if err != nil {
		return nil, err
	}

	return publisher, nil

}
