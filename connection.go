package mongomock

import (
	"sync"
)

// TestConnection is the struct that implements db.Connection
type TestConnection struct {
	m           sync.Mutex
	collections map[string]*Collection
}

// NewDB returns a testing database connection that is compatible with db.Connection
func NewDB() *TestConnection {
	return &TestConnection{
		collections: map[string]*Collection{},
	}
}

func (c *TestConnection) Collection(name string) *Collection {
	c.m.Lock()
	defer c.m.Unlock()

	v, ok := c.collections[name]
	if ok {
		return v
	}

	newCollection := &Collection{
		name:                  name,
		underlayingCollection: c,
		documents:             []documentT{},
	}
	c.collections[name] = newCollection
	return newCollection
}
