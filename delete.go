package mongomock

import (
	"github.com/mjarkk/mongomock/match"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DeleteByID deletes a document by it's ID
func (c *Collection) DeleteByID(ids ...primitive.ObjectID) error {
	c.m.Lock()
	defer c.m.Unlock()

	query := bson.M{"_id": bson.M{"$in": ids}}
	for idx := len(c.documents) - 1; idx >= 0; idx-- {
		document := c.documents[idx]
		if match.Match(document.bson, query) {
			c.documents = append(c.documents[:idx], c.documents[idx+1:]...)
			return nil
		}
	}

	return mongo.ErrNoDocuments
}
