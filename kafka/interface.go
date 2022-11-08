package kafka

import (
	"sync"

	pubsub "github.com/paper-trade-chatbot/be-pubsub"
)

type KAFKA_EXPORTED string

const (
	KAFKA_EXPORTED_NONE      KAFKA_EXPORTED = "none"
	KAFKA_EXPORTED_PRIVATE   KAFKA_EXPORTED = "private"
	KAFKA_EXPORTED_PROTECTED KAFKA_EXPORTED = "protected"
	KAFKA_EXPORTED_INTERNAL  KAFKA_EXPORTED = "internal"
	KAFKA_EXPORTED_PUBLIC    KAFKA_EXPORTED = "public"
)

type KAFKA_DATA_TYPE string

const (
	KAFKA_DATA_TYPE_NONE     KAFKA_DATA_TYPE = "none"
	KAFKA_DATA_TYPE_JSON     KAFKA_DATA_TYPE = "json"
	KAFKA_DATA_TYPE_TEXT     KAFKA_DATA_TYPE = "text"
	KAFKA_DATA_TYPE_PROTOBUF KAFKA_DATA_TYPE = "protobuf"
	KAFKA_DATA_TYPE_CSV      KAFKA_DATA_TYPE = "csv"
	KAFKA_DATA_TYPE_LOG      KAFKA_DATA_TYPE = "log"
)

type KAFKA_CONSUMER_GROUP_TYPE int

const (
	KAFKA_CONSUMER_GROUP_TYPE_NONE KAFKA_CONSUMER_GROUP_TYPE = iota
	KAFKA_CONSUMER_GROUP_TYPE_ID
	KAFKA_CONSUMER_GROUP_TYPE_REDIS
)

type Config interface {
	GetBrokers() []string
	GetExported() KAFKA_EXPORTED
	GetDomain() []string
	GetDataset() string
	GetDataType() KAFKA_DATA_TYPE
	GetVersion() uint
}

type Pubsub interface {
	pubsub.Pubsub
	GetTopicName() string
	GetMutex() *sync.Mutex
	ListTopics() []string
}

type Publisher interface {
	Pubsub
	pubsub.Publisher
}

type Subscriber interface {
	Pubsub
	pubsub.Subscriber
}
