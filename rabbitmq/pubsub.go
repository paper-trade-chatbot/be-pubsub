package rabbitmq

import (
	"context"
	"fmt"
	"sync"

	pubsub "github.com/paper-trade-chatbot/be-pubsub"
	"github.com/streadway/amqp"
)

func NewPubsub(config Config, this Pubsub) (*PubsubImpl, error) {
	pubsub := &PubsubImpl{
		Config:  config,
		Pubsub:  this,
		LogMode: pubsub.LogMode,
	}

	url := "amqp://" + config.getUsername() + ":" + config.getPassord() + "@" + config.GetHost()
	if config.GetVirtualHost() != "" {
		url += "/" + config.GetVirtualHost()
	}
	connection, err := amqp.Dial(url)
	if err != nil {
		pubsub.Log("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	channel, err := connection.Channel()
	if err != nil {
		pubsub.Log("Failed to open a channel: %v", err)
		connection.Close()
		return nil, err
	}

	pubsub.Connection = connection
	pubsub.Channel = channel

	ctx, cancel := context.WithCancel(context.Background())
	pubsub.Context = ctx
	pubsub.CancelFunc = cancel

	return pubsub, nil
}

type PubsubImpl struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Config     Config
	Pubsub
	LogMode    bool
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func (p *PubsubImpl) GetMutex() *sync.Mutex {
	return p.Pubsub.GetMutex()
}

func (p *PubsubImpl) SetLogMode(on bool) {
	p.LogMode = on
}

func (p *PubsubImpl) Log(format string, a ...any) {
	if p.LogMode {
		fmt.Printf("rabbitmq: "+format, a...)
		fmt.Println()
	}
}

func (p *PubsubImpl) Close() error {
	p.CancelFunc()
	p.Channel.Close()
	p.Connection.Close()
	return nil
}
