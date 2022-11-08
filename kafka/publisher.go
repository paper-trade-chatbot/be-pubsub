package kafka

import (
	"context"

	pubsub "github.com/paper-trade-chatbot/be-pubsub"
	"github.com/segmentio/kafka-go"
	"golang.org/x/exp/slices"
)

type PublisherImpl struct {
	PubsubImpl
	Writer *kafka.Writer
}

func NewPublisher(config *PublisherConfig, this Pubsub) (*PublisherImpl, error) {
	publisher := &PublisherImpl{
		PubsubImpl: PubsubImpl{
			Config: config,
		},
	}
	publisher.Pubsub = this

	topics := publisher.ListTopics()
	name := publisher.GetTopicName()
	publisher.Log("topics on kafka: %+v", topics)
	publisher.Log("current topic name: %+v", name)

	if !slices.Contains(topics, name) {
		publisher.Log("topic not exist in kafka: %+v", name)
		return nil, pubsub.TopicNotExist + pubsub.PubsubError("["+name+"]")
	}

	publisher.Writer = &kafka.Writer{
		Addr:         kafka.TCP(config.GetBrokers()...),
		Topic:        publisher.GetTopicName(),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: config.BatchTimeout,
	}

	publisher.Log("publisher created!")
	return publisher, nil
}

func (p *PublisherImpl) Produce(ctx context.Context, model interface{}, args ...interface{}) (interface{}, error) {
	publisher, ok := p.PubsubImpl.Pubsub.(Publisher)
	if !ok {
		return nil, pubsub.TypeCastingFailed
	}
	return publisher.Produce(ctx, model, args)
}

func (p *PublisherImpl) Close() error {
	return p.Writer.Close()
}
