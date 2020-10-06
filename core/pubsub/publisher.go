package pubsub

import (
	"fmt"
	"reflect"
	"sync"
)

type Publisher interface {
	Has(topic string) bool
	Publish(topic string, args ...interface{})
	PublishAs(topicValue ...interface{})
	Subscribe(topic string, handler interface{}) Subscriber
	SubscribeAs(handler interface{}) Subscriber
}

func NewPublisher() Publisher {
	return &publisherImpl{
		topics: make(map[string][]*subscriberImpl),
		mu:     sync.Mutex{},
	}
}

func WrapPublisher(p Publisher) Publisher {
	return &publisherImpl{
		parent: p,
		topics: make(map[string][]*subscriberImpl),
		mu:     sync.Mutex{},
	}
}

type publisherImpl struct {
	parent Publisher
	topics map[string][]*subscriberImpl
	mu     sync.Mutex
}

func (p *publisherImpl) Has(topic string) bool {
	subs, found := p.topics[topic]

	if !found || len(subs) == 0 {
		return p.parent.Has(topic)
	}

	return true
}

func (p *publisherImpl) Publish(topic string, args ...interface{}) {
	if p.parent != nil {
		p.parent.Publish(topic, args...)
	}

	if subs, found := p.topics[topic]; found && len(subs) > 0 {

		// Create reflected args
		valueArgs := make([]reflect.Value, len(args))
		for idx, arg := range args {
			valueArgs[idx] = reflect.ValueOf(arg)
		}

		for _, sub := range subs {
			if sub.handler.Type().NumIn() != len(valueArgs) {
				continue
			}

			sub.handler.Call(valueArgs)
		}
	}
}

func (p *publisherImpl) PublishAs(topicValue ...interface{}) {
	p.Publish(reflect.TypeOf(topicValue[0]).String(), topicValue...)
}

func (p *publisherImpl) Subscribe(topic string, handler interface{}) Subscriber {
	p.mu.Lock()
	defer p.mu.Unlock()

	handlerVal := reflect.ValueOf(handler)
	if handlerVal.Kind() != reflect.Func {
		panic(fmt.Errorf("cannot subscrible with %v as handler, must be a function", handlerVal.Kind()))
	}

	sub := &subscriberImpl{
		topic:     topic,
		publisher: p,
		handler:   handlerVal,
	}

	p.topics[topic] = append(p.topics[topic], sub)

	return sub
}

func (p *publisherImpl) SubscribeAs(handler interface{}) Subscriber {
	handlerType := reflect.TypeOf(handler)

	if handlerType.Kind() != reflect.Func {
		panic(fmt.Errorf("cannot subscrible with %v as handler, must be a function", handlerType.Kind()))
	}

	if handlerType.NumIn() == 0 {
		panic(fmt.Errorf("cannot subscribe with %v as handler, must have at least 1 input parameter", handlerType))
	}

	return p.Subscribe(handlerType.In(0).String(), handler)
}
