package mongomock

import (
	"github.com/mjarkk/mongomock/match"
	"github.com/mjarkk/mongomock/update"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *Collection) UpdateFirst(filter bson.M, updateValue bson.M) error {
	if updateValue == nil {
		return nil
	}

	c.m.Lock()
	defer c.m.Unlock()

	// Find document
	foundIdx := -1
	for idx, document := range c.documents {
		if match.Match(document.bson, filter) {
			foundIdx = idx
			break
		}
	}
	if foundIdx == -1 {
		return mongo.ErrNoDocuments
	}

	// Apply update to document
	document := c.documents[foundIdx]
	newBson, err := update.ApplyUpdate(document.bson, updateValue)
	if err != nil {
		return err
	}
	bsonBytes, err := bson.Marshal(document.bson)
	if err != nil {
		return err
	}
	c.documents[foundIdx] = documentT{
		bson:  newBson,
		bytes: bsonBytes,
	}

	return nil
}
