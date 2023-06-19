package mongomock

import (
	"github.com/mjarkk/mongomock/match"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ReplaceFirst updates the first document in the database that matches the filter
func (c *Collection) ReplaceFirst(filter bson.M, value any) error {
	c.m.Lock()
	defer c.m.Unlock()

	replacementDocument, err := tryNewDocument(value)
	if err != nil {
		return err
	}

	for i, entry := range c.documents {
		if match.Match(entry.bson, filter) {
			c.documents[i] = replacementDocument
			return nil
		}
	}

	return mongo.ErrNoDocuments
}

// ReplaceFirstByID updates a document in the database by its ID
// The query used here is {"_id": id}
func (c *Collection) ReplaceFirstByID(id primitive.ObjectID, value any) error {
	return c.ReplaceFirst(bson.M{"_id": id}, value)
}
