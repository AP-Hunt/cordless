package messagebus

import (
	"reflect"

	vardius_messagebus "github.com/vardius/message-bus"
)

type MessageBus interface {
	vardius_messagebus.MessageBus
	SubscribeOnce(topic string, fn interface{}) error
}

type extendedMessageBus struct {
	vardius_messagebus.MessageBus
}

func NewMessageBus() MessageBus {
	underlyingBus := vardius_messagebus.New(100)
	return &extendedMessageBus{
		MessageBus: underlyingBus,
	}
}

func (mb *extendedMessageBus) SubscribeOnce(topic string, fn interface{}) error {
	callable := reflect.ValueOf(fn)
	var wrappingCallback func(args ...interface{})

	wrappingCallback = func(args ...interface{}) {
		reflectedArgs := buildHandlerArgs(args)
		callable.Call(reflectedArgs)
		mb.Unsubscribe(topic, wrappingCallback)
	}
	return mb.Subscribe(topic, wrappingCallback)
}

// Taken from github.com/vardius/message-bus package
func buildHandlerArgs(args []interface{}) []reflect.Value {
	reflectedArgs := make([]reflect.Value, 0)

	for _, arg := range args {
		reflectedArgs = append(reflectedArgs, reflect.ValueOf(arg))
	}

	return reflectedArgs
}
