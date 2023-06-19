package match

import (
	"reflect"
	"strings"

	"github.com/mjarkk/mongomock/reflectutils"
)

// LookupMapKey looks up a key in a map.
// Note that the key can also be a nested key like "foo.bar.baz".
func lookupMapKey(scope reflect.Value, key string) *reflect.Value {
	nestedFilterkeyParts := strings.Split(key, ".")
	for _, part := range nestedFilterkeyParts {
		unwrappedScope, isNil := reflectutils.MightUnwrapPointersAndInterfaces(scope)
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
	scope, _ = reflectutils.MightUnwrapPointersAndInterfaces(scope)
	return &scope
}
