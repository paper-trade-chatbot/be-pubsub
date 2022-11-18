package rabbitmq

type ConfigImpl struct {
	Username string
	Password string

	// Host
	//
	// host path without last /
	Host string

	// VirtualHost
	//
	// namespace of mq
	// without first and last /
	VirtualHost string

	// Exported
	//
	// private: in same product
	// public: for public to access
	Exported RABBITMQ_EXPORTED

	// Domain
	//
	// domain of dataset
	Domain []string

	// Dataset
	//
	// name of this data
	Dataset string

	// ContentType
	//
	// message converted in rabbitmq
	ContentType RABBITMQ_CONTENT_TYPE

	// Version
	//
	// Version of this message
	Version uint

	// Description
	//
	// description of dataset
	Description []string
}

func (c *ConfigImpl) getUsername() string {
	return c.Username
}

func (c *ConfigImpl) getPassord() string {
	return c.Password
}

func (c *ConfigImpl) GetHost() string {
	return c.Host
}

func (c *ConfigImpl) GetVirtualHost() string {
	return c.VirtualHost
}

type PublisherConfig struct {
	ConfigImpl

	// ExchangeName
	//
	// name of this exchange
	ExchangeName string

	// ExchangeType
	//
	// type of this exchange
	ExchangeType RABBITMQ_EXCHANGE_TYPE

	// MaxRetryCount
	//
	// fail after retry over this count
	MaxRetryCount int
}

func (PublisherConfig) isConfig() {}

type SubscriberConfig struct {
	ConfigImpl

	// Consumer
	//
	// tag of this consumer
	Consumer string
}

func (SubscriberConfig) isConfig() {}
