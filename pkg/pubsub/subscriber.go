package pubsub

import "reflect"

type Subscriber interface {
	Unsubscribe()
}

type subscriberImpl struct {
	topic     string
	publisher *publisherImpl

	handler reflect.Value
}

func (s *subscriberImpl) Unsubscribe() {
	subs := s.publisher.topics[s.topic]
	if subs == nil || len(subs) == 0 {
		return
	}

	for idx, sub := range subs {
		if sub == s {
			s.publisher.topics[s.topic] = append(s.publisher.topics[s.topic][:idx], s.publisher.topics[s.topic][idx+1:]...)
		}
	}
}
