package eventst

import "sync"

type EventStream[T any] struct {
	mu sync.Mutex

	lobbies map[int32]map[string]chan<- T
}

func NewEventStream[T any]() *EventStream[T] {
	return &EventStream[T]{
		lobbies: map[int32]map[string]chan<- T{},
	}
}

func (es *EventStream[T]) Subscribe(id int32, subscriberID string, ch chan<- T) {
	es.mu.Lock()
	defer es.mu.Unlock()

	if _, ok := es.lobbies[id]; !ok {
		es.lobbies[id] = make(map[string]chan<- T, 10)
	}

	es.lobbies[id][subscriberID] = ch
}

func (es *EventStream[T]) Unsubscribe(id int32, subscriberID string) {
	es.mu.Lock()
	defer es.mu.Unlock()

	close(es.lobbies[id][subscriberID])
	delete(es.lobbies[id], subscriberID)

	if len(es.lobbies[id]) == 0 {
		delete(es.lobbies, id)
	}
}

func (es *EventStream[T]) SendEvent(id int32, d T) {
	es.mu.Lock()
	defer es.mu.Unlock()

	lobby, ok := es.lobbies[id]
	if !ok {
		return
	}

	for _, subscriber := range lobby {
		subscriber <- d
	}
}
