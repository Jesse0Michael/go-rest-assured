package assured

import (
	"sync"
)

type CallStore struct {
	data map[string][]*Call
	sync.Mutex
}

func NewCallStore() *CallStore {
	return &CallStore{data: map[string][]*Call{}}
}

func (c *CallStore) Add(call *Call) {
	c.Lock()
	c.data[call.ID()] = append(c.data[call.ID()], call)
	c.Unlock()
}

func (c *CallStore) AddAt(key string, call *Call) {
	c.Lock()
	c.data[key] = append(c.data[key], call)
	c.Unlock()
}

func (c *CallStore) Rotate(call *Call) {
	c.Lock()
	c.data[call.ID()] = append(c.data[call.ID()][1:], call)
	c.Unlock()
}

func (c *CallStore) Get(key string) []*Call {
	c.Lock()
	calls := c.data[key]
	c.Unlock()
	return calls
}

func (c *CallStore) Clear(key string) {
	c.Lock()
	delete(c.data, key)
	c.Unlock()
}

func (c *CallStore) ClearAll() {
	c.Lock()
	c.data = map[string][]*Call{}
	c.Unlock()
}
