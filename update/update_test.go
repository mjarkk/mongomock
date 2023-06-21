package update

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func assertApplyUpdate(t *testing.T, doucment, update, expectedResults bson.M) {
	actualResults, err := ApplyUpdate(doucment, update)
	NoError(t, err)
	Equal(t, expectedResults, actualResults)
}

func TestSetUpdate(t *testing.T) {
	// Add a property
	assertApplyUpdate(
		t,
		bson.M{},
		bson.M{"$set": bson.M{"bar": 20}},
		bson.M{"bar": 20},
	)

	// Update property
	assertApplyUpdate(
		t,
		bson.M{"bar": 10},
		bson.M{"$set": bson.M{"bar": 20}},
		bson.M{"bar": 20},
	)

	// Update nested property using nested keys
	assertApplyUpdate(
		t,
		bson.M{"foo": bson.M{"bar": 10}},
		bson.M{"$set": bson.M{"foo.bar": 20}},
		bson.M{"foo": bson.M{"bar": 20}},
	)

	// Update nested property using a map
	assertApplyUpdate(
		t,
		bson.M{"foo": bson.M{"bar": 10}},
		bson.M{"$set": bson.M{"foo": bson.M{"bar": 20}}},
		bson.M{"foo": bson.M{"bar": 20}},
	)

	// Update nested property that does not exist
	assertApplyUpdate(
		t,
		bson.M{},
		bson.M{"$set": bson.M{"foo.bar.baz": 20}},
		bson.M{"foo": bson.M{"bar": bson.M{"baz": 20}}},
	)

	// Set nested properties using numeric indexes
	assertApplyUpdate(
		t,
		bson.M{},
		bson.M{"$set": bson.M{"0": 20}},
		bson.M{"0": 20},
	)

	// Update array index
	assertApplyUpdate(
		t,
		bson.M{"foo": []int{1, 2, 3}},
		bson.M{"$set": bson.M{"foo.1": 20}},
		bson.M{"foo": []int{1, 20, 3}},
	)

	// Update array index to another type
	assertApplyUpdate(
		t,
		bson.M{"foo": []int{1, 2, 3}},
		bson.M{"$set": bson.M{"foo.1": "bogus"}},
		bson.M{"foo": []any{1, "bogus", 3}},
	)
}
