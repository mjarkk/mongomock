package reflectutils

import (
	"fmt"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

// SliceLikeToSlice converts a slice-like value to a slice of any.
// A slice-like value is either a slice or an array or a pointer to a slice or an array.
// Nil slice-like values are returned as empty slices.
func SliceLikeToSlice(sliceLike any) (slice []any, isSliceLike bool) {
	sliceReflection, isNil := MightUnwrapPointersAndInterfaces(reflect.ValueOf(sliceLike))
	kind := sliceReflection.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, false
	}
	if isNil || (kind == reflect.Slice && sliceReflection.IsNil()) {
		return []any{}, true
	}

	slice = make([]any, sliceReflection.Len())
	for idx := 0; idx < sliceReflection.Len(); idx++ {
		slice[idx] = sliceReflection.Index(idx).Interface()
	}

	return slice, true
}

// MightUnwrapPointersAndInterfaces tries to unwrap pointers from a value.
func MightUnwrapPointersAndInterfaces(v reflect.Value) (unwrappedValue reflect.Value, isNil bool) {
outer:
	for {
		switch v.Kind() {
		case reflect.Ptr:
			if v.IsNil() {
				isNil = true
			}
			v = v.Elem()
		case reflect.Interface:
			v = v.Elem()
		default:
			break outer
		}
	}

	return v, isNil
}

// ConvertMapToBsonM converts a map structure to bson.M
// `m` Must be of kind reflect.Map
func ConvertMapToBsonM(m reflect.Value) (bson.M, error) {
	result := bson.M{}
	mRange := m.MapRange()
	for mRange.Next() {
		k, err := NormalizeKey(mRange.Key())
		if err != nil {
			return nil, err
		}
		result[k] = mRange.Value().Interface()
	}
	return result, nil
}

// NormalizeKey converts a key like to a string.
func NormalizeKey(key reflect.Value) (string, error) {
	key, _ = MightUnwrapPointersAndInterfaces(key)

	switch key.Kind() {
	case reflect.String:
		return key.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.Itoa(int(key.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.Itoa(int(key.Uint())), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(key.Float(), 'f', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(key.Bool()), nil
	default:
		return "", fmt.Errorf("key must be a string or number, got %s", key.Type())
	}
}
