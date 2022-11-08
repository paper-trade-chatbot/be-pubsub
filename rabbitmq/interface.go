package rabbitmq

import (
	"sync"

	pubsub "github.com/paper-trade-chatbot/be-pubsub"
)

type RABBITMQ_EXPORTED string

const (
	RABBITMQ_EXPORTED_NONE     RABBITMQ_EXPORTED = "none"
	RABBITMQ_EXPORTED_PRIVATE  RABBITMQ_EXPORTED = "private"
	RABBITMQ_EXPORTED_INTERNAL RABBITMQ_EXPORTED = "internal"
	RABBITMQ_EXPORTED_PUBLIC   RABBITMQ_EXPORTED = "public"
)

type RABBITMQ_EXCHANGE_TYPE string

const (
	RABBITMQ_EXCHANGE_TYPE_NONE    RABBITMQ_EXCHANGE_TYPE = "none"
	RABBITMQ_EXCHANGE_TYPE_DIRECT  RABBITMQ_EXCHANGE_TYPE = "direct"
	RABBITMQ_EXCHANGE_TYPE_FANOUT  RABBITMQ_EXCHANGE_TYPE = "fanout"
	RABBITMQ_EXCHANGE_TYPE_TOPIC   RABBITMQ_EXCHANGE_TYPE = "topic"
	RABBITMQ_EXCHANGE_TYPE_HEADERS RABBITMQ_EXCHANGE_TYPE = "headers"
)

type RABBITMQ_CONTENT_TYPE string

const (
	RABBITMQ_CONTENT_TYPE_NONE RABBITMQ_CONTENT_TYPE = ""
	RABBITMQ_CONTENT_TYPE_JSON RABBITMQ_CONTENT_TYPE = "application/json"
	RABBITMQ_CONTENT_TYPE_TEXT RABBITMQ_CONTENT_TYPE = "text/plain"
)

type Config interface {
	isConfig()
	getUsername() string
	getPassord() string
	GetHost() string
	GetVirtualHost() string
}

type Pubsub interface {
	pubsub.Pubsub
	GetMutex() *sync.Mutex
}

type Publisher[T interface{}] interface {
	Pubsub
	pubsub.TPublisher[T]
	GetExchangeName() string
	GetExchangeType() RABBITMQ_EXCHANGE_TYPE
	GetRoutingKey() string
}

type Subscriber[T interface{}] interface {
	Pubsub
	pubsub.TSubscriber[T]
	GetQueueName() string
	GetConsumer() string
}
