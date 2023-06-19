package update

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mjarkk/mongomock/reflectutils"
	"go.mongodb.org/mongo-driver/bson"
)

func ApplyUpdate(document bson.M, update bson.M) (bson.M, error) {
	if update == nil {
		return document, nil
	}

	document, err := applySet(document, update["$set"])
	if err != nil {
		return nil, err
	}

	return document, errors.New("not implemented")
}

func applySet(document bson.M, value any) (bson.M, error) {
	valueReflection, isNil := reflectutils.MightUnwrapPointersAndInterfaces(reflect.ValueOf(value))
	if isNil {
		return document, nil
	}

	if valueReflection.Kind() != reflect.Map {
		return nil, fmt.Errorf("$set expected map, got %s", valueReflection.Kind())
	}

	if valueReflection.IsNil() {
		return document, nil
	}

	mapRange := valueReflection.MapRange()
	for mapRange.Next() {
		keyReflection := mapRange.Key()
		valueReflection := mapRange.Value()
		if keyReflection.Kind() != reflect.String {
			return nil, errors.New("$set map keys must be a string")
		}
		key := keyReflection.String()

		var err error
		document, err = bsonSetField(document, key, valueReflection)
		if err != nil {
			return nil, err
		}
	}

	return document, nil
}

func bsonSetNestedFieldBase(document bson.M, baseKeyAndSubKey []string, value reflect.Value) (bson.M, error) {
	var err error
	document[baseKeyAndSubKey[0]], err = bsonSetField(bson.M{}, baseKeyAndSubKey[1], value)
	return document, err
}

func bsonSetField(document bson.M, key string, value reflect.Value) (bson.M, error) {
	keyParts := strings.SplitN(key, ".", 2)
	if len(keyParts) == 1 {
		document[key] = value
		return document, nil
	}

	// This is a nested key (foo.bar.baz)
	// Note that nested key might also target array elements (foo.0.baz)
	field, ok := document[keyParts[0]]
	fieldReflection, isNil := reflectutils.MightUnwrapPointersAndInterfaces(reflect.ValueOf(field))
	if !ok || isNil {
		return bsonSetNestedFieldBase(document, keyParts, value)
	}

	switch fieldReflection.Kind() {
	case reflect.Map:
		typedField, ok := field.(bson.M)
		if !ok {
			return nil, fmt.Errorf("trying to access a nested key on a non-bson.M map: %s", keyParts[0])
		}
		return bsonSetField(typedField, keyParts[1], value)
	case reflect.Slice, reflect.Array:
		value, err := bsonSetSliceLike(fieldReflection, key, value)
		if err != nil {
			return nil, err
		}
		document[keyParts[0]] = value
		return document, nil
	default:
		return nil, errors.New("trying to access a nested key on a non-map, non-slice, non-array field")
	}
}

func bsonSetSliceLike(potentialDocumentSlice reflect.Value, key string, value reflect.Value) (any, error) {
	keyParts := strings.SplitN(key, ".", 2)
	isNestedKey := len(keyParts) == 2

	switch potentialDocumentSlice.Kind() {
	case reflect.Slice:
		if potentialDocumentSlice.IsNil() {
			// Slice is nil (null in mongodb) thus we can try to reassign it
			if isNestedKey {
				return nil, fmt.Errorf("trying to set a nested property on a null value, key: %s", key)
			}
			return value.Interface(), nil
		}
		fallthrough
	case reflect.Array:
		// Assign to the index of the array
		arrayIndex, err := strconv.Atoi(keyParts[0])
		if err != nil {
			return nil, fmt.Errorf("trying to access a slice like element using a non-integer key: %s", key)
		}
		if arrayIndex < 0 {
			return nil, fmt.Errorf("trying to access a slice like element using a negative key: %s", key)
		}
		if arrayIndex >= potentialDocumentSlice.Len() {
			// The array is too small, we need to resize it
			if potentialDocumentSlice.Kind() == reflect.Array {
				return nil, fmt.Errorf("trying to access a slice like element using an out of range key (This should be possible but the document contains a staticly sized array): %s", keyParts[0])
			}
			potentialDocumentSlice.SetLen(arrayIndex + 1)
		}

		valueToSet := value
		if isNestedKey {
			nestedValue, err := bsonSetField(bson.M{}, keyParts[1], value)
			if err != nil {
				return nil, err
			}
			valueToSet = reflect.ValueOf(nestedValue)
		}
		potentialDocumentSlice.Index(arrayIndex).Set(valueToSet)

		return potentialDocumentSlice.Interface(), nil
	default:
		return nil, fmt.Errorf("trying to access a slice like element on a non-slice like field: %s", key)
	}
}
