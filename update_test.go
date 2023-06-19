package mongomock

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUpdate(t *testing.T) {
	testingDB := NewDB()
	usersCollection := testingDB.Collection("users")

	document := bson.M{
		"foo": "bar",
		"bar": 10,
	}

	err := usersCollection.Insert(document)
	NoError(t, err)

	err = usersCollection.UpdateFirst(bson.M{"foo": "bar"}, bson.M{
		"$set": bson.M{
			"bar": 20,
		},
	})
	NoError(t, err)
}
