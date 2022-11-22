package pubsub

import "context"

var LogMode bool = false

type Pubsub interface {
	SetLogMode(bool)
	Log(string, ...any)
	Close() error
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
	Pubsub
	Produce(context.Context, T, ...interface{}) (interface{}, error)
}

type TSubscriber[T interface{}] interface {
	Pubsub
	Listen(context.Context, ...interface{}) error
	Consume(context.Context, T, ...interface{}) (interface{}, error)
	Subscribe(context.Context, func(context.Context, T) error) error
}
