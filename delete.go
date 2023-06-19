package mongomock

import (
	"github.com/mjarkk/mongomock/match"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DeleteFirst deletes the first document that matches the filter
func (c *Collection) DeleteFirst(filter bson.M) error {
	c.m.Lock()
	defer c.m.Unlock()

	for idx, document := range c.documents {
		if match.Match(document.bson, filter) {
			c.documents = append(c.documents[:idx], c.documents[idx+1:]...)
			return nil
		}
	}

	return mongo.ErrNoDocuments
}

// Delete deletes all documents matching the filter
// The query used here is {"_id": {"$in": ids}}
func (c *Collection) Delete(filter bson.M) error {
	c.m.Lock()
	defer c.m.Unlock()

	for idx := len(c.documents) - 1; idx >= 0; idx-- {
		document := c.documents[idx]
		if match.Match(document.bson, filter) {
			c.documents = append(c.documents[:idx], c.documents[idx+1:]...)
			return nil
		}
	}

	return mongo.ErrNoDocuments
}

// DeleteByID deletes a document by it's ID
// The query used here is {"_id": id}
func (c *Collection) DeleteByID(id primitive.ObjectID) error {
	return c.Delete(bson.M{"_id": id})
}

// DeleteByIDs deletes documents by their IDs
// The query used here is {"_id": {"$in": ids}}
func (c *Collection) DeleteByIDs(ids ...primitive.ObjectID) error {
	return c.Delete(bson.M{"_id": bson.M{"$in": ids}})
}
