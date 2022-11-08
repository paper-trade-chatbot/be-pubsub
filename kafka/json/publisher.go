package json

import (
	"context"
	"encoding/json"

	"github.com/asaskevich/govalidator"
	pubsub "github.com/paper-trade-chatbot/be-pubsub"
	pubsubKafka "github.com/paper-trade-chatbot/be-pubsub/kafka"
	"github.com/segmentio/kafka-go"
)

type PublisherImpl struct {
	pubsubKafka.PublisherImpl
}

func NewPublisher(config *pubsubKafka.PublisherConfig) (*PublisherImpl, error) {

	if config.GetDataType() != pubsubKafka.KAFKA_DATA_TYPE_JSON {
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
// model must be a struct, otherwise it won't work
func (p *PublisherImpl) Produce(ctx context.Context, model interface{}, args ...interface{}) (interface{}, error) {
	var err error

	if _, err := govalidator.ValidateStruct(model); err != nil {
		return nil, err
	}

	message := []byte{}
	if message, err = json.Marshal(model); err != nil {
		return nil, err
	}
	p.Log("[" + p.GetTopicName() + "] producing: " + string(message))

	kMessage := kafka.Message{
		Value: message,
	}
	err = p.PublisherImpl.Writer.WriteMessages(ctx, kMessage)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
