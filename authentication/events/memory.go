package events

import (
	"context"
	"sync"
)

type MemoryBus struct {
	mu          sync.RWMutex
	nextID      uint64
	subscribers map[EventName]map[uint64]Handler
}

func NewMemoryBus() *MemoryBus {
	return &MemoryBus{subscribers: make(map[EventName]map[uint64]Handler)}
}

func (b *MemoryBus) Publish(ctx context.Context, event Envelope) error {
	if b == nil {
		return nil
	}

	b.mu.RLock()
	handlers := make([]Handler, 0, len(b.subscribers[event.Name]))
	for _, handler := range b.subscribers[event.Name] {
		handlers = append(handlers, handler)
	}
	b.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

func (b *MemoryBus) Subscribe(name EventName, handler Handler) Subscription {
	if b == nil || handler == nil {
		return noopSubscription{}
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.nextID++
	id := b.nextID
	if b.subscribers == nil {
		b.subscribers = make(map[EventName]map[uint64]Handler)
	}
	if b.subscribers[name] == nil {
		b.subscribers[name] = make(map[uint64]Handler)
	}
	b.subscribers[name][id] = handler

	return &memorySubscription{bus: b, name: name, id: id}
}

type memorySubscription struct {
	bus  *MemoryBus
	name EventName
	id   uint64
	once sync.Once
}

func (s *memorySubscription) Unsubscribe() {
	s.once.Do(func() {
		s.bus.mu.Lock()
		defer s.bus.mu.Unlock()

		delete(s.bus.subscribers[s.name], s.id)
		if len(s.bus.subscribers[s.name]) == 0 {
			delete(s.bus.subscribers, s.name)
		}
	})
}
