package match

import (
	"reflect"
	"strings"
)

// sliceLikeToSlice converts a slice-like value to a slice of any.
// A slice-like value is either a slice or an array or a pointer to a slice or an array.
// Nil slice-like values are returned as empty slices.
func sliceLikeToSlice(sliceLike any) (slice []any, isSliceLike bool) {
	sliceReflection, isNil := mightUnwrapPointersAndInterfaces(reflect.ValueOf(sliceLike))
	if sliceReflection.Kind() != reflect.Slice && sliceReflection.Kind() != reflect.Array {
		return nil, false
	}
	if isNil || sliceReflection.IsNil() {
		return []any{}, true
	}

	slice = make([]any, sliceReflection.Len())
	for idx := 0; idx < sliceReflection.Len(); idx++ {
		slice[idx] = sliceReflection.Index(idx).Interface()
	}

	return slice, true
}

// mightUnwrapPointersAndInterfaces tries to unwrap pointers from a value.
func mightUnwrapPointersAndInterfaces(v reflect.Value) (unwrappedValue reflect.Value, isNil bool) {
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

// lookupMapKey looks up a key in a map.
// Note that the key can also be a nested key like "foo.bar.baz".
func lookupMapKey(scope reflect.Value, key string) *reflect.Value {
	nestedFilterkeyParts := strings.Split(key, ".")
	for _, part := range nestedFilterkeyParts {
		unwrappedScope, isNil := mightUnwrapPointersAndInterfaces(scope)
		if isNil {
			return nil
		}

		switch unwrappedScope.Kind() {
		case reflect.Struct:
			unwrappedScope = reflect.ValueOf(mustConvertToBson(unwrappedScope.Interface()))
			if unwrappedScope.IsNil() {
				return nil
			}
		case reflect.Map:
			if unwrappedScope.IsNil() {
				return nil
			}
			// continue
		default:
			return nil
		}

		scope = unwrappedScope.MapIndex(reflect.ValueOf(part))
		if !scope.IsValid() {
			return nil
		}
	}
	scope, _ = mightUnwrapPointersAndInterfaces(scope)
	return &scope
}
