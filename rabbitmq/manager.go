package rabbitmq

import (
	"reflect"
	"sync"
)

type manager struct {
	Publishers  sync.Map
	Subscribers sync.Map
}

var Manager manager

func GetPublisher[T any]() Publisher[T] {

	var t T
	if p, ok := Manager.Publishers.Load(reflect.TypeOf(t).String()); ok {
		if _, ok = p.(Publisher[T]); ok {
			return p.(Publisher[T])
		}
	}

	return nil
}

func SetPublisher[T any](publisher Publisher[T]) bool {

	var t T
	Manager.Publishers.Store(reflect.TypeOf(t).String(), publisher)

	return true
}
