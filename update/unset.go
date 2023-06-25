package update

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mjarkk/mongomock/reflectutils"
	"go.mongodb.org/mongo-driver/bson"
)

func applyUnset(document bson.M, update any) (bson.M, error) {
	if update == nil {
		return document, nil
	}

	// FIXME Support update as string or map of keys

	updateReflection := reflect.ValueOf(update)
	updateReflection, isNil := reflectutils.MightUnwrapPointersAndInterfaces(updateReflection)
	if isNil {
		return document, nil
	}

	switch updateReflection.Kind() {
	case reflect.String:
		return unsetKeys(document, []string{updateReflection.String()})
	case reflect.Map:
		if updateReflection.IsNil() {
			return document, nil
		}

		mapRange := updateReflection.MapRange()
		keys := []string{}
		for mapRange.Next() {
			key, err := reflectutils.NormalizeKey(mapRange.Key())
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)
		}

		return unsetKeys(document, keys)
	case reflect.Slice:
		if updateReflection.IsNil() {
			return document, nil
		}
		fallthrough
	case reflect.Array:
		keys := []string{}
		for idx := 0; idx < updateReflection.Len(); idx++ {
			keyReflection := updateReflection.Index(idx)
			key, err := reflectutils.NormalizeKey(keyReflection)
			if err != nil {
				return nil, fmt.Errorf("%s, index: %d", err.Error(), idx)
			}
			keys = append(keys, key)
		}

		return unsetKeys(document, keys)
	default:
		return nil, fmt.Errorf("$unset expected a slice like, got %s", updateReflection.Type())
	}
}

func unsetKeys(document bson.M, keys []string) (bson.M, error) {
	var err error
	for _, key := range keys {
		document, err = unsetKey(document, key)
		if err != nil {
			return nil, fmt.Errorf("%s, key: %s", err.Error(), key)
		}
	}
	return document, nil
}

func unsetKey(document bson.M, key string) (bson.M, error) {
	keyParts := strings.SplitN(key, ".", 2)

	if len(keyParts) == 1 {
		delete(document, key)
		return document, nil
	}

	entry, ok := document[keyParts[0]]
	if !ok {
		// Nothing to delete
		return document, nil
	}

	entryReflection, isNil := reflectutils.MightUnwrapPointersAndInterfaces(reflect.ValueOf(entry))
	if isNil {
		// Nothing to delete
		return document, nil
	}

	switch entryReflection.Kind() {
	case reflect.Map:
		if entryReflection.IsNil() {
			// Nothing to delete
			return document, nil
		}

		childDocument, err := reflectutils.ConvertMapToBsonM(entryReflection)
		if err != nil {
			return nil, err
		}

		result, err := unsetKey(childDocument, keyParts[1])
		if err != nil {
			return nil, err
		}

		document[keyParts[0]] = result
		return document, nil
	case reflect.Slice:
		if entryReflection.IsNil() {
			// Nothing to delete
			return document, nil
		}
		fallthrough
	case reflect.Array:
		err := unsetSliceLike(&entryReflection, keyParts[1])
		document[keyParts[0]] = entryReflection.Interface()
		return document, err
	default:
		// This is a value without any subfields, so we can't delete anything
		return document, nil
	}
}

// unsetSliceLike unsets a key in a slice like (slice or array)
// There are 2 paths here
// 1. The key is a nested key, so we need to try recurse into the nested document
// 2. This is the end of a nested key, so we need to delete this key for the slice
func unsetSliceLike(sliceLike *reflect.Value, key string) error {
	keyParts := strings.SplitN(key, ".", 2)

	index, err := strconv.Atoi(keyParts[0])
	if err != nil {
		return err
	}

	if index >= sliceLike.Len() {
		return nil
	}

	if len(keyParts) == 1 {
		// Path 2
		sliceLikeElementType := sliceLike.Type().Elem()
		switch sliceLikeElementType.Kind() {
		case reflect.Ptr | reflect.Interface | reflect.Map | reflect.Slice:
			// null value is actually null (nil)
		default:
			// null value is actually the zero value
			*sliceLike = convertSliceLikeToInterfaceSlice(*sliceLike)
		}

		sliceLike.Index(index).SetZero()
		return nil
	}

	// Path 1
	entry := sliceLike.Index(index)
	unwrappedEntry, isNil := reflectutils.MightUnwrapPointersAndInterfaces(entry)
	if isNil || entry.Kind() != reflect.Map || entry.IsNil() {
		return fmt.Errorf("expected a map as this is a nested key %s, got %v", key, entry.Interface())
	}

	entryAsMap, err := reflectutils.ConvertMapToBsonM(unwrappedEntry)
	if err != nil {
		return err
	}

	entryAsMap, err = unsetKey(entryAsMap, keyParts[1])
	if err != nil {
		return err
	}

	entry.Set(reflect.ValueOf(entryAsMap))

	return nil
}
