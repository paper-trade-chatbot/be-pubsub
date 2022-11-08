package log

import (
	"context"

	pubsub "github.com/paper-trade-chatbot/be-pubsub"
	pubsubKafka "github.com/paper-trade-chatbot/be-pubsub/kafka"
	"github.com/segmentio/kafka-go"
)

type PublisherImpl struct {
	pubsubKafka.PublisherImpl
}

func NewPublisher(config *pubsubKafka.PublisherConfig) (*PublisherImpl, error) {

	if config.GetDataType() != pubsubKafka.KAFKA_DATA_TYPE_LOG {
		return nil, pubsub.WrongDataType
	}

	publisher := &PublisherImpl{}

	if parent, err := pubsubKafka.NewPublisher(config, publisher); err != nil {
		publisher.Log("failed to new publisher: %v", err)
		return nil, err
	} else {
		publisher.PublisherImpl = *parent
	}

	return publisher, nil
}

// Produce
//
// model must be a string, otherwise it won't work
func (p *PublisherImpl) Produce(ctx context.Context, model interface{}, args ...interface{}) (interface{}, error) {
	var err error

	modelString, ok := model.(string)
	if !ok {
		p.Log("Produce: input model is not pointer to string. %v")
		return nil, pubsub.InputNotString
	}

	p.Log("[" + p.GetTopicName() + "] producing: " + modelString)

	kMessage := kafka.Message{
		Value: []byte(modelString),
	}
	err = p.PublisherImpl.Writer.WriteMessages(ctx, kMessage)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
