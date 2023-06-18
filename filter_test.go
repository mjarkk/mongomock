package mongomock

import (
	"testing"
	"time"

	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FilterMatches(filter bson.M, data any) bool {
	return newFilter(filter).matches(data)
}

func TestFilter(t *testing.T) {
	stringValue := "abc"

	type exampleNestedField struct {
		Bar string
	}

	scenarios := []struct {
		name              string
		matchingFilter    bson.M
		nonMatchingFilter bson.M
		data              any
	}{
		{
			"empty filter",
			bson.M{},
			bson.M{"a": true},
			struct{}{},
		},
		{
			"bool field match",
			bson.M{"foo": true},
			bson.M{"foo": false},
			struct{ Foo bool }{true},
		},
		{
			"int field match",
			bson.M{"foo": 123},
			bson.M{"foo": 1},
			struct{ Foo int16 }{123},
		},
		{
			"string field match",
			bson.M{"foo": "123"},
			bson.M{"foo": "abc"},
			struct{ Foo string }{"123"},
		},
		{
			"bson tag",
			bson.M{"bar": "123"},
			bson.M{"foo": "123"},
			struct {
				Foo string `bson:"bar"`
			}{"123"},
		},
		{
			"pointer value",
			bson.M{"foo": "abc"},
			bson.M{"foo": nil},
			struct {
				Foo *string
			}{&stringValue},
		},
		{
			"object id",
			bson.M{"foo": primitive.ObjectID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}},
			bson.M{"foo": primitive.ObjectID{}},
			struct {
				Foo primitive.ObjectID
			}{primitive.ObjectID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}},
		},
		{
			"inline test",
			bson.M{"bar": "abc"},
			bson.M{"foo": "abc"},
			struct {
				Foo exampleNestedField `bson:",inline"`
			}{exampleNestedField{"abc"}},
		},
		{
			"$gt with int",
			bson.M{"foo": bson.M{"$gt": 5}},
			bson.M{"foo": bson.M{"$gt": 10}},
			struct{ Foo int }{Foo: 7},
		},
		{
			"$lt with int",
			bson.M{"foo": bson.M{"$lt": 10}},
			bson.M{"foo": bson.M{"$lt": 5}},
			struct{ Foo int }{Foo: 7},
		},
		{
			"$gt with uint",
			bson.M{"foo": bson.M{"$gt": 5}},
			bson.M{"foo": bson.M{"$gt": 10}},
			struct{ Foo uint }{Foo: 7},
		},
		{
			"$lt with uint",
			bson.M{"foo": bson.M{"$lt": 10}},
			bson.M{"foo": bson.M{"$lt": 5}},
			struct{ Foo uint }{Foo: 7},
		},
		{
			"$gte",
			bson.M{"foo": bson.M{"$gte": 7}},
			bson.M{"foo": bson.M{"$gte": 10}},
			struct{ Foo int }{Foo: 7},
		},
		{
			"$lte",
			bson.M{"foo": bson.M{"$lte": 7}},
			bson.M{"foo": bson.M{"$lte": 5}},
			struct{ Foo int }{Foo: 7},
		},
		{
			"$gt with time",
			bson.M{"foo": bson.M{"$gt": time.Now()}},
			bson.M{"foo": bson.M{"$gt": time.Now().Add(time.Hour)}},
			struct{ Foo time.Time }{Foo: time.Now().Add(time.Minute * 30)},
		},
		{
			"$lt with time",
			bson.M{"foo": bson.M{"$lt": time.Now().Add(time.Hour)}},
			bson.M{"foo": bson.M{"$lt": time.Now()}},
			struct{ Foo time.Time }{Foo: time.Now().Add(time.Minute * 30)},
		},
		{
			"$eq",
			bson.M{"foo": bson.M{"$eq": 5}},
			bson.M{"foo": bson.M{"$eq": 2}},
			struct{ Foo int }{Foo: 5},
		},
		{
			"$or",
			bson.M{"$or": []bson.M{
				{"foo": 1},
				{"foo": 2},
				{"foo": 3},
			}},
			bson.M{"$or": []bson.M{
				{"foo": 1},
				{"foo": 3},
			}},
			struct{ Foo int }{Foo: 2},
		},
		{
			"$and",
			bson.M{"$and": []bson.M{
				{"foo": 2},
				{"bar": 1},
			}},
			bson.M{"$and": []bson.M{
				{"foo": 2},
				{"bar": 0},
			}},
			struct {
				Foo int
				Bar int
			}{Foo: 2, Bar: 1},
		},
		{
			"$not",
			bson.M{"foo": bson.M{"$not": 1}},
			bson.M{"foo": bson.M{"$not": 2}},
			struct{ Foo int }{Foo: 2},
		},
		{
			"$size",
			bson.M{"foo": bson.M{"$size": 2}},
			bson.M{"foo": bson.M{"$size": 3}},
			struct{ Foo []int }{Foo: []int{1, 2}},
		},
		{
			"$bitsAllSet",
			bson.M{"foo": bson.M{"$bitsAllSet": uint(1)}},
			bson.M{"foo": bson.M{"$bitsAllSet": uint(1) << 1}},
			struct{ Foo uint }{Foo: 1},
		},
		{
			"$type",
			bson.M{"foo": bson.M{"$type": 2}},
			bson.M{"foo": bson.M{"$type": 3}},
			struct{ Foo string }{},
		},
		{
			"path with dots",
			bson.M{"foo.bar": "abc"},
			bson.M{"foo.bas": "abc"},
			struct{ Foo exampleNestedField }{Foo: exampleNestedField{Bar: "abc"}},
		},
		{
			"filter slice of numbers",
			bson.M{"foo": 2},
			bson.M{"foo": 4},
			struct{ Foo []int }{Foo: []int{1, 2, 3}},
		},
		{
			"filter slice of numbers using a specific index",
			bson.M{"foo.1": 2},
			bson.M{"foo.2": 2},
			struct{ Foo []int }{Foo: []int{1, 2, 3}},
		},
		{
			"filter slice with objects",
			bson.M{"foo": bson.M{"bar": "bbb"}},
			bson.M{"foo": bson.M{"bar": "zzz"}},
			struct{ Foo []exampleNestedField }{[]exampleNestedField{{"aaa"}, {"bbb"}, {"ccc"}}},
		},
		{
			"filter slice with objects using dotted paths",
			bson.M{"foo.bar": "bbb"},
			bson.M{"foo.bar": "zzz"},
			struct{ Foo []exampleNestedField }{[]exampleNestedField{{"aaa"}, {"bbb"}, {"ccc"}}},
		},
	}

	for idx := range scenarios {
		t.Run(scenarios[idx].name, func(t *testing.T) {
			s := scenarios[idx]
			True(t, FilterMatches(s.matchingFilter, s.data))
			False(t, FilterMatches(s.nonMatchingFilter, s.data))
		})
	}
}

func TestFilterType(t *testing.T) {
	// From:
	// https://docs.mongodb.com/manual/reference/operator/query/type/

	scenarios := []struct {
		typeID      int
		typeName    string
		dataToMatch any
	}{
		{1, "double", 1.0},
		{2, "string", "foo"},
		{3, "object", struct{}{}},
		{4, "array", []int{1, 2}},
		// 5 binData
		// 6 undefined
		// 7 objectId
		{8, "bool", true},
		// 9 date
		{10, "null", nil},
		// 11 regex
		// 12 dbPointer
		// 13 javascript
		// 14 symbol
		// 15 javascriptWithScope
		{16, "int", 1},
		// 17 timestamp
		// 18 long
		{19, "decimal", 1.0},
		// -1 minKey
		// 127 maxKey
	}

	for _, s := range scenarios {
		t.Run(s.typeName, func(t *testing.T) {
			True(t, FilterMatches(bson.M{"foo": bson.M{"$type": s.typeName}}, struct{ Foo any }{Foo: s.dataToMatch}))
			True(t, FilterMatches(bson.M{"foo": bson.M{"$type": s.typeID}}, struct{ Foo any }{Foo: s.dataToMatch}))
		})
	}
}
