package kafka

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type ConfigImpl struct {
	// Brokers
	//
	// Address of the kafka cluster that this writer is configured to send
	// messages to.
	//
	// This field is required, attempting to write messages to a writer with a
	// nil address will error.
	Brokers []string

	// Exported
	//
	// private: in same product
	// public: for public to access
	Exported KAFKA_EXPORTED

	// Domain
	//
	// domain of dataset
	Domain []string

	// Dataset
	//
	// name of this data
	Dataset string

	// DataType
	//
	// message converted in kafka
	DataType KAFKA_DATA_TYPE

	// Version
	//
	// Version of this topic
	Version uint
}

func (c *ConfigImpl) GetBrokers() []string {
	return c.Brokers
}

func (c *ConfigImpl) GetExported() KAFKA_EXPORTED {
	return c.Exported
}

func (c *ConfigImpl) GetDomain() []string {
	return c.Domain
}

func (c *ConfigImpl) GetDataset() string {
	return c.Dataset
}

func (c *ConfigImpl) GetDataType() KAFKA_DATA_TYPE {
	return c.DataType
}

func (c *ConfigImpl) GetVersion() uint {
	return c.Version
}

type PublisherConfig struct {
	Config

	// Time limit on how often incomplete message batches will be flushed to
	// kafka.
	//
	// The default is to flush at least every second.
	// minimum: 10ms. if less than 10ms, this won't work
	BatchTimeout time.Duration
}

type SubscriberConfig struct {
	Config
	// ConsumerGroupType
	//
	// how to maintain offset and make sure only access once in a group.
	// we could either use "kafka consumer group" way or "redis" way
	ConsumerGroupType KAFKA_CONSUMER_GROUP_TYPE

	// Redis
	//
	// if we use redis to maintain consumer group and offset, this should be specified
	Redis *redis.Client

	// GroupID holds the optional consumer group id.  If GroupID is specified, then
	// Partition should NOT be specified e.g. 0
	ConsumerGroupID string

	// RetentionTime optionally sets the length of time the consumer group will be saved
	// by the broker
	//
	// Default: 24h
	//
	// Only used when GroupID is set
	RetentionTime time.Duration

	// Maximum amount of time to wait for new data to come when fetching batches
	// of messages from kafka.
	//
	// Default: 10s
	// minimum: 10ms. if less than 10ms, this won't work
	MaxWait time.Duration

	// MinBytes indicates to the broker the minimum batch size that the consumer
	// will accept. Setting a high minimum when consuming from a low-volume topic
	// may result in delayed delivery when the broker does not have enough data to
	// satisfy the defined minimum.
	//
	// Default: 1
	MinBytes int

	// MaxBytes indicates to the broker the maximum batch size that the consumer
	// will accept. The broker will truncate a message to satisfy this maximum, so
	// choose a value that is high enough for your largest message size.
	//
	// Default: 1MB
	MaxBytes int
}
