package rabbitmq

import (
	"github.com/paper-trade-chatbot/be-pubsub/rabbitmq"
	"github.com/paper-trade-chatbot/be-pubsub/rabbitmq/json"
)

func NewPublisher(username, password, host, virualHost string) (*json.PublisherImpl[*NewPositionModel], error) {

	publisher, err := json.NewPublisher[*NewPositionModel](
		&rabbitmq.PublisherConfig{
			ConfigImpl: rabbitmq.ConfigImpl{
				Username:    username,
				Password:    password,
				Host:        host,
				VirtualHost: virualHost,
				Exported:    rabbitmq.RABBITMQ_EXPORTED_PRIVATE,
				Domain:      []string{"order"},
				Dataset:     "newPosition",
				ContentType: rabbitmq.RABBITMQ_CONTENT_TYPE_JSON,
			},
			ExchangeName: "order",
			ExchangeType: rabbitmq.RABBITMQ_EXCHANGE_TYPE_DIRECT,
		})

	if err != nil {
		return nil, err
	}

	return publisher, nil

}