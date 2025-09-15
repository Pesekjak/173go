package evt

import (
	"sync"
)

type EventManager struct {
	idCounter int
	listeners map[string][]RegisteredListener
	mutex     sync.RWMutex
}

type Event struct {
	Cancelled bool
	Type      string
	Data      interface{}
}

type EventListener func(event *Event)

type RegisteredListener struct {
	id int
	EventListener
	Unsubscribe func() bool
}

func NewEventManager() *EventManager {
	return &EventManager{
		idCounter: 0,
		listeners: make(map[string][]RegisteredListener),
	}
}

func (em *EventManager) AddListener(eventType string, listener EventListener) RegisteredListener {
	em.mutex.Lock()
	defer em.mutex.Unlock()
	registeredListener := RegisteredListener{
		id:            em.idCounter,
		EventListener: listener,
		Unsubscribe:   nil,
	}
	registeredListener.Unsubscribe = func() bool {
		em.mutex.Lock()
		defer em.mutex.Unlock()
		if registeredListener.id == -1 {
			return false
		}
		listeners, ok := em.listeners[eventType]
		if !ok {
			return false
		}
		for i, nextListener := range listeners {
			if nextListener.id != registeredListener.id {
				continue
			}
			em.listeners[eventType] = append(listeners[:i], listeners[i+1:]...)
			registeredListener.id = -1
			return true
		}
		return false
	}
	em.listeners[eventType] = append(em.listeners[eventType], registeredListener)
	em.idCounter++
	return registeredListener
}

func (em *EventManager) FireEvent(event Event) bool {
	eventType := event.Type
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	listeners, ok := em.listeners[eventType]
	if !ok {
		return true // no listeners fired, no event was cancelled
	}

	for _, registered := range listeners {
		registered.EventListener(&event)
		if event.Cancelled {
			return false
		}
	}

	return true
}
