package rabbitmq

import (
	"strconv"
	"strings"

	"github.com/streadway/amqp"
)

type PublisherImpl struct {
	PubsubImpl
	Jobs chan *PublishJob
}

type PublishJob struct {
	Message    *amqp.Publishing
	RetryCount int
}

func NewPublisher[T interface{}](config *PublisherConfig, this Publisher[T]) (*PublisherImpl, error) {
	pubsub, err := NewPubsub(config, this)
	if err != nil {
		return nil, err
	}

	publisher := &PublisherImpl{
		PubsubImpl: *pubsub,
		Jobs:       make(chan *PublishJob, 128),
	}

	if _, err = publisher.Channel.QueueDeclare(
		publisher.GetRoutingKey(), // name
		false,                     // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		nil,                       // arguments
	); err != nil {
		publisher.Log("channel[%s] declare failed", publisher.GetRoutingKey())
		return nil, err
	}

	publisher.Log("publisher created! [%s][%s]", publisher.GetExchangeName(), publisher.GetRoutingKey())

	go func() {
		for {
			select {
			case j := <-publisher.Jobs:
				publisher.Log("publish [%s][%s] message: %s", publisher.GetExchangeName(), publisher.GetRoutingKey(), string(j.Message.Body))
				err = publisher.Channel.Publish(
					publisher.GetExchangeName(), // exchange
					publisher.GetRoutingKey(),   // routing key
					false,                       // mandatory
					false,                       // immediate
					*j.Message)
				if err != nil {
					publisher.Log("Failed to publish [%s][%s]. retry count %d", publisher.GetExchangeName(), publisher.GetRoutingKey(), j.RetryCount)
					j.RetryCount++
					publisher.Jobs <- j
					continue
				}
				publisher.Log("[" + publisher.GetExchangeName() + "][" + publisher.GetRoutingKey() + "] producing: " + string(j.Message.Body))
			case <-publisher.Context.Done():
				publisher.Log("publisher [%s][%s] terminated.", publisher.GetExchangeName(), publisher.GetRoutingKey())
				return
			}
		}
	}()

	return publisher, nil
}

func (p *PublisherImpl) GetExchangeName() string {
	config := p.Config.(*PublisherConfig)
	return config.ExchangeName + "." + string(config.ExchangeType)
}

func (p *PublisherImpl) GetExchangeType() RABBITMQ_EXCHANGE_TYPE {
	config := p.Config.(*PublisherConfig)
	return config.ExchangeType
}

func (p *PublisherImpl) GetRoutingKey() string {

	config := p.Config.(*PublisherConfig)

	name := []string{}

	name = append(name, string(config.Exported))

	name = append(name, config.Domain...)
	if len(config.Domain) == 0 {
		name = append(name, "default")
	}

	name = append(name, config.Dataset)

	name = append(name, strconv.FormatUint(uint64(config.Version), 10))

	name = append(name, config.Description...)

	keyName := strings.Join(name, ".")

	return keyName
}
