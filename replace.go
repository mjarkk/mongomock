package mongomock

import (
	"github.com/mjarkk/mongomock/match"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateByID updates a document in the database by its ID
func (c *Collection) ReplaceOneById(id primitive.ObjectID, value any) error {
	c.m.Lock()
	defer c.m.Unlock()

	replacementDocument, err := TryNewDocument(value)
	if err != nil {
		return err
	}

	query := bson.M{"_id": id}

	for i, entry := range c.documents {
		if match.Match(entry, query) {
			c.documents[i] = replacementDocument
			break
		}
	}

	return nil
}
