package config

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
)

const DBEventTopic = "db_config"
const LogEventTopic = "log_config"
const APPEventTopic = "app_config"

type Event struct {
	Topic string
	Val   interface{}
}

type Observer interface {
	OnChange() error
}

type EventBus struct {
	mux       sync.RWMutex
	observers map[string]map[Observer]struct{}
}

type BaseObserver struct {
}

func (b *BaseObserver) OnChange(e *Event) error {
	log.Printf("observer: %s, event key: %s, event val: %v", b, e.Topic, e.Val)
	return nil
}

var eventBus = &EventBus{
	observers: make(map[string]map[Observer]struct{}),
}

func GetEventBus() *EventBus {
	return eventBus
}

func (s *EventBus) Subscribe(topic string, o Observer) {
	s.mux.Lock()
	defer s.mux.Unlock()
	_, ok := s.observers[topic]
	if !ok {
		s.observers[topic] = make(map[Observer]struct{})
	}
	s.observers[topic][o] = struct{}{}
}

func (s *EventBus) Unsubscribe(topic string, o Observer) {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.observers[topic], o)
}

func (s *EventBus) Publish(e *Event) error {
	s.mux.RLock()
	defer s.mux.RUnlock()
	subscribers := s.observers[e.Topic]

	errs := make(map[Observer]error)
	for subscriber := range subscribers {
		if err := subscriber.OnChange(); err != nil {
			errs[subscriber] = err
		}
	}

	return s.handleErr(errs)
}

func (s *EventBus) handleErr(errs map[Observer]error) error {
	if len(errs) > 0 {
		message := ""
		for o, err := range errs {
			message += fmt.Sprintf("observer: %v, err: %v;", o, err)
		}

		return errors.New(message)
	}

	return nil
}

func PublishEvents(newConfig Config, topics []string) error {
	errMsg := ""
	for _, topic := range topics {
		err := eventBus.Publish(&Event{
			Topic: topic,
			Val:   newConfig,
		})
		if err != nil {
			errMsg += err.Error()
		}
	}

	if errMsg != "" {
		return errors.New(errMsg)
	}

	return nil
}
