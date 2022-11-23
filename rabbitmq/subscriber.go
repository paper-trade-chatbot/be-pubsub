package rabbitmq

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/streadway/amqp"
)

type SubscriberImpl[T interface{}] struct {
	PubsubImpl
	ListenMutex *sync.RWMutex
	Delivery    <-chan amqp.Delivery
	Callbacks   []func(context.Context, T) error
}

func NewSubscriber[T interface{}](config *SubscriberConfig, this Subscriber[T]) (*SubscriberImpl[T], error) {
	pubsub, err := NewPubsub(config, this)
	if err != nil {
		return nil, err
	}

	subscriber := &SubscriberImpl[T]{
		PubsubImpl:  *pubsub,
		ListenMutex: &sync.RWMutex{},
	}

	if _, err := pubsub.Channel.QueueDeclare(
		subscriber.GetQueueName(), // name
		true,                      // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		nil,                       // arguments
	); err != nil {
		subscriber.Log("queue not exist in rabbitmq: %s", subscriber.GetQueueName())
		return nil, err
	}

	delivery, err := pubsub.Channel.Consume(
		subscriber.GetQueueName(), // queue
		subscriber.GetConsumer(),  // consumer
		true,                      // auto-ack
		false,                     // exclusive
		false,                     // no-local
		false,                     // no-wait
		nil,                       // args
	)
	if err != nil {
		pubsub.Close()
		return nil, err
	}

	subscriber.Delivery = delivery

	subscriber.Log("subscriber created! [%s][%s]", subscriber.GetQueueName(), subscriber.GetConsumer())
	return subscriber, nil
}

func (s *SubscriberImpl[T]) Subscribe(ctx context.Context, callback func(context.Context, T) error) error {
	s.Callbacks = append(s.Callbacks, callback)
	return nil
}

func (s *SubscriberImpl[T]) GetQueueName() string {
	config := s.Config.(*SubscriberConfig)

	name := []string{}

	name = append(name, string(config.Exported))

	name = append(name, config.Domain...)
	if len(config.Domain) == 0 {
		name = append(name, "default")
	}

	name = append(name, config.Dataset)

	name = append(name, strconv.FormatUint(uint64(config.Version), 10))

	name = append(name, config.Description...)

	queueName := strings.Join(name, ".")

	return queueName
}

func (s *SubscriberImpl[T]) GetConsumer() string {
	config := s.Config.(*SubscriberConfig)
	return config.Consumer
}
