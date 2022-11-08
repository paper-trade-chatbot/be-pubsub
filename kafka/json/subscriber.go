package json

import (
	"context"
	"encoding/json"
	"reflect"
	"runtime/debug"

	"github.com/asaskevich/govalidator"
	pubsub "github.com/paper-trade-chatbot/be-pubsub"
	pubsubKafka "github.com/paper-trade-chatbot/be-pubsub/kafka"
)

type SubscriberImpl struct {
	pubsubKafka.SubscriberImpl
}

func NewSubscriber(config *pubsubKafka.SubscriberConfig) (*SubscriberImpl, error) {
	if config.GetDataType() != pubsubKafka.KAFKA_DATA_TYPE_JSON {
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

	// if model is a struct, then turn into pointer
	if reflect.ValueOf(model).Type().Kind() != reflect.Pointer {
		s.Log("error: input model must be pointer of struct.")
		return pubsub.ListenNotPointer
	}

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

					s.Log("message: %s", string(m.Value))

					// clear the model
					p := reflect.ValueOf(model).Elem()
					p.Set(reflect.Zero(p.Type()))

					for _, f := range s.Callbacks {

						// clone a model
						clone := reflect.New(p.Type()).Interface()
						if err := json.Unmarshal(m.Value, clone); err != nil {
							s.Log("error: failed to unmarshal. %v", err)
							continue
						}

						if _, err := govalidator.ValidateStruct(reflect.ValueOf(clone)); err != nil {
							s.Log("error: ValidateStruct err:%v\n", err)
							continue
						}

						if err := f(clone); err != nil {
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

	if reflect.ValueOf(model).Type().Kind() != reflect.Pointer {
		s.Log("error: input model must be pointer of struct.")
		return nil, pubsub.ListenNotPointer
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

	if err := json.Unmarshal(m.Value, model); err != nil {
		s.Log("Consume: failed to unmarshal. %v", err)
		return nil, err
	}

	if _, err := govalidator.ValidateStruct(model); err != nil {
		s.Log("Consume: ValidateStruct err:%v\n", err)
		return nil, err
	}

	return nil, nil
}
