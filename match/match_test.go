package match

import (
	"fmt"
	"testing"

	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func mustObjectIDFromHex(hex string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		panic(err.Error())
	}
	return id
}

func TestMatch(t *testing.T) {
	cases := []struct {
		Name              string
		Document          bson.M
		MatchingFilter    bson.M
		NotMatchingFilter bson.M
	}{
		{
			"simple",
			bson.M{"foo": "bar"},
			bson.M{"foo": "bar"},
			bson.M{"bar": "foo"},
		},
		{
			"simple with multiple fields",
			bson.M{"foo": "bar", "bar": "baz"},
			bson.M{"foo": "bar"},
			bson.M{"foo": "bar", "bar": "foo"},
		},
		{
			"nested query 1",
			bson.M{"foo": bson.M{"bar": "baz"}},
			bson.M{"foo.bar": "baz"},
			bson.M{"foo.bar": "foo"},
		},
		{
			"nested query 2",
			bson.M{"foo": bson.M{"bar": "baz"}},
			bson.M{"foo": bson.M{"bar": "baz"}},
			bson.M{"foo": bson.M{"bar": "foo"}},
		},
		{
			"array contains",
			bson.M{"foo": []string{"bar", "baz"}},
			bson.M{"foo": "bar"},
			bson.M{"foo": "foo"},
		},
		{
			"array matches",
			bson.M{"foo": []string{"bar", "baz"}},
			bson.M{"foo": []string{"bar", "baz"}},
			bson.M{"foo": []string{"bar", "foo"}},
		},
		{
			"object ids",
			bson.M{"_id": mustObjectIDFromHex("5a0e7f6c0f291d0f1f2d3e4f")},
			bson.M{"_id": mustObjectIDFromHex("5a0e7f6c0f291d0f1f2d3e4f")},
			bson.M{"_id": mustObjectIDFromHex("000000000000000000000000")},
		},
		{
			"$eq",
			bson.M{"foo": "bar"},
			bson.M{"foo": bson.M{"$eq": "bar"}},
			bson.M{"foo": bson.M{"$eq": "baz"}},
		},
		{
			"$ne",
			bson.M{"foo": "bar"},
			bson.M{"foo": bson.M{"$ne": "baz"}},
			bson.M{"foo": bson.M{"$ne": "bar"}},
		},
		{
			"$gt",
			bson.M{"foo": 4},
			bson.M{"foo": bson.M{"$gt": 2}},
			bson.M{"foo": bson.M{"$gt": 4}},
		},
		{
			"$gte",
			bson.M{"foo": 4},
			bson.M{"foo": bson.M{"$gte": 4}},
			bson.M{"foo": bson.M{"$gte": 5}},
		},
		{
			"$lt",
			bson.M{"foo": 4},
			bson.M{"foo": bson.M{"$lt": 6}},
			bson.M{"foo": bson.M{"$lt": 4}},
		},
		{
			"$lte",
			bson.M{"foo": 4},
			bson.M{"foo": bson.M{"$lte": 4}},
			bson.M{"foo": bson.M{"$lte": 3}},
		},
		{
			"$and",
			bson.M{"foo": 4, "bar": 5},
			bson.M{"$and": []bson.M{{"foo": 4}, {"bar": 5}}},
			bson.M{"$and": []bson.M{{"foo": 4}, {"bar": 6}}},
		},
		{
			"$size",
			bson.M{"foo": []string{"bar", "baz"}},
			bson.M{"foo": bson.M{"$size": 2}},
			bson.M{"foo": bson.M{"$size": 3}},
		},
		{
			"$type",
			bson.M{"foo": "bar"},
			bson.M{"foo": bson.M{"$type": "string"}},
			bson.M{"foo": bson.M{"$type": "int"}},
		},
		{
			"$all",
			bson.M{"foo": []string{"foo", "bar", "baz"}},
			bson.M{"foo": bson.M{"$all": []string{"foo", "baz"}}},
			bson.M{"foo": bson.M{"$all": []string{"foo", "bar", "baz", "qux"}}},
		},
		{
			"$all with custom data types",
			bson.M{"foo": []any{"foo", 2, 1.0}},
			bson.M{"foo": bson.M{"$all": []any{"foo", 1.0}}},
			bson.M{"foo": bson.M{"$all": []any{"foo", 3}}},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.Name+" matches case", func(t *testing.T) {
			match := internalMatch(testCase.Document, testCase.MatchingFilter)

			hint := fmt.Sprintf("filter %+v should match document %+v", testCase.MatchingFilter, testCase.Document)
			True(t, match, hint)
		})
		t.Run(testCase.Name+" not matches case", func(t *testing.T) {
			match := internalMatch(testCase.Document, testCase.NotMatchingFilter)

			hint := fmt.Sprintf("filter %+v should NOT match document %+v", testCase.NotMatchingFilter, testCase.Document)
			False(t, match, hint)
		})
	}
}
