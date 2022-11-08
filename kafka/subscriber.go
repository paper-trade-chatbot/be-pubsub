package kafka

import (
	"context"
	"sync"

	pubsub "github.com/paper-trade-chatbot/be-pubsub"
	"github.com/segmentio/kafka-go"
	"golang.org/x/exp/slices"
)

type SubscriberImpl struct {
	PubsubImpl
	Reader      *kafka.Reader
	ListenMutex *sync.RWMutex
	Callbacks   []func(interface{}) error
}

func NewSubscriber(config *SubscriberConfig, this Pubsub) (*SubscriberImpl, error) {
	subscriber := &SubscriberImpl{
		PubsubImpl: PubsubImpl{
			Config: config,
		},
		ListenMutex: &sync.RWMutex{},
	}
	subscriber.Pubsub = this

	topics := subscriber.ListTopics()
	name := subscriber.GetTopicName()
	subscriber.Log("topics on kafka: %+v", topics)
	subscriber.Log("current topic name: %+v", name)

	if !slices.Contains(topics, name) {
		subscriber.Log("topic not exist in kafka: %+v", name)
		return nil, pubsub.TopicNotExist + pubsub.PubsubError("["+name+"]")
	}

	subscriber.Reader = kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:  config.GetBrokers(),
			Topic:    subscriber.GetTopicName(),
			MinBytes: config.MaxBytes,
			MaxBytes: config.MaxBytes,
			MaxWait:  config.MaxWait,
			GroupID:  config.ConsumerGroupID,
		},
	)

	subscriber.Log("subscriber created!")
	return subscriber, nil
}

// Listen
//
// model must be a pointer to a struct, otherwise it won't work
func (s *SubscriberImpl) Listen(ctx context.Context, model interface{}, args ...interface{}) error {
	subscriber, ok := s.PubsubImpl.Pubsub.(Subscriber)
	if !ok {
		return pubsub.TypeCastingFailed
	}
	return subscriber.Listen(ctx, model, args)
}

func (s *SubscriberImpl) Consume(ctx context.Context, model interface{}, args ...interface{}) (interface{}, error) {
	subscriber, ok := s.PubsubImpl.Pubsub.(Subscriber)
	if !ok {
		return nil, pubsub.TypeCastingFailed
	}
	return subscriber.Consume(ctx, model, args)
}

func (s *SubscriberImpl) Subscribe(ctx context.Context, callback func(interface{}) error) error {
	s.Callbacks = append(s.Callbacks, callback)
	return nil
}

func (s *SubscriberImpl) Close() error {
	return s.Reader.Close()
}
