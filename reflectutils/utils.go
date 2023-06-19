package reflectutils

import (
	"reflect"
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
