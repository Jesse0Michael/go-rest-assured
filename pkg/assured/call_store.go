package assured

import (
	"sync"
)

type Selector interface {
	Key() string
}

type Store[T Selector] struct {
	data map[string][]T
	sync.Mutex
}

func NewStore[T Selector]() *Store[T] {
	return &Store[T]{data: map[string][]T{}}
}

func (c *Store[T]) Add(v T) {
	c.Lock()
	c.data[v.Key()] = append(c.data[v.Key()], v)
	c.Unlock()
}

func (c *Store[T]) AddAt(key string, call T) {
	c.Lock()
	c.data[key] = append(c.data[key], call)
	c.Unlock()
}

func (c *Store[T]) Rotate(v T) {
	c.Lock()
	c.data[v.Key()] = append(c.data[v.Key()][1:], v)
	c.Unlock()
}

func (c *Store[T]) Get(key string) []T {
	c.Lock()
	calls := c.data[key]
	c.Unlock()
	return calls
}

func (c *Store[T]) Clear(key string) {
	c.Lock()
	delete(c.data, key)
	c.Unlock()
}

func (c *Store[T]) ClearAll() {
	c.Lock()
	c.data = map[string][]T{}
	c.Unlock()
}
