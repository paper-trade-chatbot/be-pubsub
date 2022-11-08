package rabbitmq

import (
	"context"

	rfPubsub "github.com/paper-trade-chatbot/be-pubsub"
	"github.com/paper-trade-chatbot/be-pubsub/rabbitmq"
	"github.com/paper-trade-chatbot/be-pubsub/rabbitmq/json"
)

func NewSubscriber(username, password, host, virualHost, consumer string) (*json.SubscriberImpl[*NewPositionModel], error) {
	subscriber, err := json.NewSubscriber[*NewPositionModel](
		&rabbitmq.SubscriberConfig{
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
			Consumer: consumer,
		})

	if err != nil {
		return nil, err
	}

	return subscriber, nil
}

// SubscribeAndListen
//
// model must be a pointer to a struct, otherwise it won't work
func SubscribeAndListen(ctx context.Context, username, password, host, virualHost, consumer string, callbacks ...func(*NewPositionModel) error) (rfPubsub.TSubscriber[*NewPositionModel], error) {
	if len(callbacks) == 0 {
		return nil, rfPubsub.ListenNullCallback
	}

	sub, err := NewSubscriber(username, password, host, virualHost, consumer)
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

	err = sub.Listen(ctx)
	if err != nil {
		sub.Close()
		return nil, err
	}
	return sub, nil
}
