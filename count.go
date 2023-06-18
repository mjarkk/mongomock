package mongomock

import (
	"github.com/mjarkk/mongomock/match"
	"go.mongodb.org/mongo-driver/bson"
)

// Count returns the number of documents in the collection of entity
func (c *Collection) Count(filter bson.M) (uint64, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if len(filter) == 0 {
		// Take the easy route
		return uint64(len(c.documents)), nil
	}

	var count uint64
	for _, document := range c.documents {
		if match.Match(document, filter) {
			count++
		}
	}

	return count, nil
}
