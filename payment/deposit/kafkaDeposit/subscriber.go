package kafkaDeposit

import (
	"context"
	"time"

	rfPubsub "github.com/paper-trade-chatbot/be-pubsub"
	"github.com/paper-trade-chatbot/be-pubsub/kafka"
	"github.com/paper-trade-chatbot/be-pubsub/kafka/json"
)

func NewSubscriber(consumerGroup string, brokers []string) (*json.SubscriberImpl, error) {
	subscriber, err := json.NewSubscriber(&kafka.SubscriberConfig{
		Config: &kafka.ConfigImpl{
			Brokers:  brokers,
			Exported: kafka.KAFKA_EXPORTED_PRIVATE,
			Domain:   []string{"payment"},
			Dataset:  "deposit",
			DataType: kafka.KAFKA_DATA_TYPE_JSON,
			Version:  0,
		},
		ConsumerGroupType: kafka.KAFKA_CONSUMER_GROUP_TYPE_ID,
		Redis:             nil,
		ConsumerGroupID:   consumerGroup,
		MaxWait:           time.Millisecond * 10,
		MinBytes:          1,
		MaxBytes:          10e6, //1MB
	})

	if err != nil {
		return nil, err
	}

	return subscriber, nil
}

// SubscribeAndListen
//
// model must be a pointer to a struct, otherwise it won't work
func SubscribeAndListen(ctx context.Context, consumerGroup string, brokers []string, model interface{}, callbacks ...func(interface{}) error) (rfPubsub.Subscriber, error) {
	if len(callbacks) == 0 {
		return nil, rfPubsub.ListenNullCallback
	}

	sub, err := NewSubscriber(
		consumerGroup,
		brokers,
	)
	if err != nil {
		return nil, err
	}

	for _, c := range callbacks {
		err = sub.Subscribe(ctx, c)
		if err != nil {
			sub.Close()
			return nil, err
		}
	}

	err = sub.Listen(ctx, model)
	if err != nil {
		sub.Close()
		return nil, err
	}
	return sub, nil
}
