package pubsub

import "context"

type Pubsub interface {
	SetLogMode(bool)
	Log(string, ...any)
}

type Publisher interface {
	Produce(context.Context, interface{}, ...interface{}) (interface{}, error)
	Close() error
}

type Subscriber interface {
	Listen(context.Context, interface{}, ...interface{}) error
	Consume(context.Context, interface{}, ...interface{}) (interface{}, error)
	Subscribe(context.Context, func(interface{}) error) error
	Close() error
}

type TPublisher[T interface{}] interface {
	Produce(context.Context, T, ...interface{}) (interface{}, error)
	Close() error
}

type TSubscriber[T interface{}] interface {
	Listen(context.Context, ...interface{}) error
	Consume(context.Context, T, ...interface{}) (interface{}, error)
	Subscribe(context.Context, func(T) error) error
	Close() error
}
