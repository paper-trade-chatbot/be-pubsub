package log

import (
	"context"
	"runtime/debug"

	pubsub "github.com/paper-trade-chatbot/be-pubsub"
	pubsubKafka "github.com/paper-trade-chatbot/be-pubsub/kafka"
)

type SubscriberImpl struct {
	pubsubKafka.SubscriberImpl
}

func NewSubscriber(config *pubsubKafka.SubscriberConfig) (*SubscriberImpl, error) {
	if config.GetDataType() != pubsubKafka.KAFKA_DATA_TYPE_LOG {
		return nil, pubsub.WrongDataType
	}

	subscriber := &SubscriberImpl{}

	if parent, err := pubsubKafka.NewSubscriber(config, subscriber); err != nil {
		subscriber.Log("failed to new subscriber: %v", err)
		return nil, err
	} else {
		subscriber.SubscriberImpl = *parent
	}

	return subscriber, nil
}

// Listen
//
// model must be a pointer to a struct, otherwise it won't work
func (s *SubscriberImpl) Listen(ctx context.Context, model interface{}, args ...interface{}) error {

	if ok := s.ListenMutex.TryLock(); !ok {
		return pubsub.ListenWhileConsuming
	}

	// TODO: cancel with context, to gracefully shut down listening

	go func() {
		defer s.ListenMutex.Unlock()

		for {
			func() {
				defer func() {
					if err := recover(); err != nil {
						s.Log("Listen panic:", err)
						s.Log("stacktrace from panic: \n" + string(debug.Stack()))
					}
				}()

				s.Log("start listening to %s...", s.GetTopicName())
				for {
					m, err := s.Reader.ReadMessage(ctx)
					if err != nil {
						s.Log("error: failed to read message. %v", err)
						continue
					}

					for _, f := range s.Callbacks {
						message := string(m.Value)

						if err := f(&message); err != nil {
							s.Log("error: Callback err:%v\n", err)
							continue
						}
					}
				}
			}()
		}

	}()

	return nil
}

func (s *SubscriberImpl) Consume(ctx context.Context, model interface{}, args ...interface{}) (interface{}, error) {

	modelString, ok := model.(*string)
	if !ok {
		s.Log("Consume: input model is not pointer to string. %v")
		return nil, pubsub.InputNotString
	}

	if ok := s.ListenMutex.TryRLock(); !ok {
		return nil, pubsub.ConsumeWhileListening
	}
	defer s.ListenMutex.RUnlock()

	m, err := s.Reader.ReadMessage(ctx)
	if err != nil {
		s.Log("Consume: failed to read message. %v", err)
		return nil, err
	}

	*modelString = string(m.Value)

	return nil, nil
}
