package mongomock

import (
	"errors"
	"reflect"

	"github.com/mjarkk/mongomock/match"
	"go.mongodb.org/mongo-driver/bson"
)

// documentT describes a document within the collection
type documentT struct {
	bson  bson.M
	bytes []byte
}

func tryNewDocument(value any) (documentT, error) {
	parsedValue, isNil := match.MightUnwrapPointersAndInterfaces(reflect.ValueOf(value))
	if isNil {
		return documentT{}, errors.New("value is nil")
	}

	encodedValue, err := bson.Marshal(parsedValue.Interface())
	if err != nil {
		return documentT{}, err
	}
	encodedValueMap := bson.M{}
	err = bson.Unmarshal(encodedValue, &encodedValueMap)
	if err != nil {
		return documentT{}, err
	}

	return documentT{
		bson:  encodedValueMap,
		bytes: encodedValue,
	}, nil
}
