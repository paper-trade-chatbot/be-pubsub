package json

import (
	"context"
	"encoding/json"

	"github.com/asaskevich/govalidator"
	pubsub "github.com/paper-trade-chatbot/be-pubsub"
	rabbitmq "github.com/paper-trade-chatbot/be-pubsub/rabbitmq"
	"github.com/streadway/amqp"
)

type PublisherImpl[T interface{}] struct {
	rabbitmq.PublisherImpl
}

func NewPublisher[T interface{}](config *rabbitmq.PublisherConfig) (*PublisherImpl[T], error) {

	if config.ContentType != rabbitmq.RABBITMQ_CONTENT_TYPE_JSON {
		return nil, pubsub.WrongDataType
	}

	publisher := &PublisherImpl[T]{}

	if parent, err := rabbitmq.NewPublisher[T](config, publisher); err != nil {
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
func (p *PublisherImpl[T]) Produce(ctx context.Context, model T, args ...interface{}) (interface{}, error) {
	var err error

	if _, err := govalidator.ValidateStruct(model); err != nil {
		return nil, err
	}

	message := []byte{}
	if message, err = json.Marshal(model); err != nil {
		return nil, err
	}

	job := &rabbitmq.PublishJob{
		Message: &amqp.Publishing{
			ContentType: string(p.Config.(*rabbitmq.PublisherConfig).ContentType),
			Body:        message,
		},
	}

	p.Jobs <- job

	return nil, nil
}
