package json

import (
	"context"
	"encoding/json"
	"reflect"
	"runtime/debug"

	"github.com/asaskevich/govalidator"
	pubsub "github.com/paper-trade-chatbot/be-pubsub"
	rabbitmq "github.com/paper-trade-chatbot/be-pubsub/rabbitmq"
)

type SubscriberImpl[T interface{}] struct {
	rabbitmq.SubscriberImpl[T]
}

func NewSubscriber[T interface{}](config *rabbitmq.SubscriberConfig) (*SubscriberImpl[T], error) {

	if config.ContentType != rabbitmq.RABBITMQ_CONTENT_TYPE_JSON {
		return nil, pubsub.WrongDataType
	}

	subscriber := &SubscriberImpl[T]{}

	if parent, err := rabbitmq.NewSubscriber[T](config, subscriber); err != nil {
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
func (s *SubscriberImpl[T]) Listen(ctx context.Context, args ...interface{}) error {

	var model T
	// if model is a struct, then turn into pointer
	if reflect.ValueOf(model).Type().Kind() != reflect.Pointer {
		s.Log("error: input model must be pointer of struct.")
		return pubsub.ListenNotPointer
	}

	if ok := s.ListenMutex.TryLock(); !ok {
		return pubsub.ListenWhileConsuming
	}

	go func() {
		defer s.ListenMutex.Unlock()
		s.Log("start listening to %s by %s...", s.GetQueueName(), s.GetConsumer())
		for {
			if func() int {
				defer func() {
					if err := recover(); err != nil {
						s.Log("Listen panic:", err)
						s.Log("stacktrace from panic: \n" + string(debug.Stack()))
					}
				}()

				select {
				case d := <-s.Delivery:

					s.Log("message: %s", string(d.Body))
					for _, f := range s.Callbacks {
						var model T
						modelType := reflect.TypeOf(model).Elem()
						model = reflect.New(modelType).Interface().(T)
						if err := json.Unmarshal(d.Body, model); err != nil {
							s.Log("error: failed to unmarshal. %v", err)
							continue
						}

						if _, err := govalidator.ValidateStruct(model); err != nil {
							s.Log("error: ValidateStruct err:%v\n", err)
							continue
						}

						if err := f(ctx, model); err != nil {
							s.Log("error: Callback err:%v\n", err)
							continue
						}

					}
				case <-s.Context.Done():
					s.Log("subscriber [%s][%s] terminated.", s.GetQueueName(), s.GetConsumer())
					return 1
				}
				return 0
			}() == 1 {
				return
			}
		}
	}()

	return nil
}

func (s *SubscriberImpl[T]) Consume(ctx context.Context, model T, args ...interface{}) (interface{}, error) {

	if reflect.ValueOf(model).Type().Kind() != reflect.Pointer {
		s.Log("error: input model must be pointer of struct.")
		return nil, pubsub.ListenNotPointer
	}

	if ok := s.ListenMutex.TryRLock(); !ok {
		return nil, pubsub.ConsumeWhileListening
	}
	defer s.ListenMutex.RUnlock()

	d, ok := <-s.Delivery
	if !ok {
		s.Log("Consume: no new message. [%v]", pubsub.NullDelivery)
		return nil, pubsub.NullDelivery
	}

	if err := json.Unmarshal(d.Body, model); err != nil {
		s.Log("error: failed to unmarshal. %v", err)
		return nil, err
	}

	if _, err := govalidator.ValidateStruct(model); err != nil {
		s.Log("Consume: ValidateStruct err:%v\n", err)
		return nil, err
	}

	return nil, nil
}
