package match

import (
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// Match matches a document against a filter
// returns true if it matches
func Match(document any, filter bson.M) bool {
	if filter == nil {
		return true
	}

	return internalMatch(document, filter)
}

func internalMatch(document any, filter bson.M) bool {
	for filterKey, filterValue := range filter {
		filterOperator, isOperator := strings.CutPrefix(filterKey, "$")
		if isOperator {
			if !valueMatchesOperator(document, filterOperator, filterValue) {
				return false
			}
			continue
		}

		documentElement := lookupMapKey(reflect.ValueOf(document), filterKey)
		var value any = nil
		if documentElement != nil {
			value = documentElement.Interface()
		}

		if !valueMatchesFilter(value, filterValue) {
			return false
		}
	}
	return true
}

// valueMatchesFilter checks if the value matches the filter
// This in the bases is just foo == bar
// But mongodb supports lots of operators and this function also resolves them
func valueMatchesFilter(value any, filter any) bool {
	valueSlice, valueIsSliceLike := sliceLikeToSlice(value)
	if valueIsSliceLike {
		matched := sliceLikeValueMatchesSliceLikeFilter(valueSlice, filter)
		if matched {
			return true
		}
		// Coninue to check if there are other ways to match
		// Like: Matching filter {$size: 3} against value [4, 2, 0]
	}

	switch typedFilter := filter.(type) {
	case bson.M:
		return internalMatch(value, typedFilter)
	case nil:
		return value == nil
	case string:
		typedValue, ok := value.(string)
		if !ok {
			return false
		}
		return typedValue == typedFilter
	case bool:
		typedValue, ok := value.(bool)
		if !ok {
			return false
		}
		return typedValue == typedFilter
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return numberValueMatchesFilter(value, filter, eqComparator)
	default:
		_, filterIsSliceLike := sliceLikeToSlice(filter)
		if filterIsSliceLike {
			return false
		}

		filterReflection, isNil := mightUnwrapPointersAndInterfaces(reflect.ValueOf(filter))
		switch filterReflection.Kind() {
		case reflect.Struct, reflect.Map:
			filter = mustConvertToBson(filter)
			if isNil {
				return valueMatchesFilter(nil, filter)
			}
			return valueMatchesFilter(value, filter)
		case reflect.Slice, reflect.Array:
			return false
		}

		panic("unknown filter type, only bson.M for objects is supported. All basic types should be supported")
	}
}

func sliceLikeValueMatchesSliceLikeFilter(valueSlice []any, filter any) bool {
	filterSlice, filterIsSliceLike := sliceLikeToSlice(filter)
	if !filterIsSliceLike {
		// Trying to match
		// Document: { age: [1, 2, 3] }
		// Query: { age: 2 }
		for _, value := range valueSlice {
			if valueMatchesFilter(value, filter) {
				return true
			}
		}
	}

	if len(filterSlice) != len(valueSlice) {
		return false
	}

	for idx := 0; idx < len(filterSlice); idx++ {
		filter := filterSlice[idx]
		value := valueSlice[idx]
		if !valueMatchesFilter(value, filter) {
			return false
		}
	}

	return true
}

// valueMatchesOperator checks if the value matches the operator filter
// Like:
//
//	Query: { age: { $gt: 5 } }
//	Document: { age: 10 }
//	Example: valueMatchesOperator(10, "$gt", 5)
//
//	Query: { age: { $in: [5, 10] } }
//	Document: { age: 10 }
//	Example: valueMatchesOperator(10, "$in", [5, 10])
//
//	Query: { $or: [ { age: 5 }, { age: 10 } ] }
//	Document: { age: 10 }
//	Example: valueMatchesOperator(10, "$or", [{age: 5}, {age: 10}])
func valueMatchesOperator(value any, operator string, operatorFilter any) bool {
	switch operator {
	case "eq":
		return valueMatchesFilter(value, operatorFilter)
	case "ne", "not":
		return !valueMatchesFilter(value, operatorFilter)
	case "gt":
		return numberValueMatchesFilter(value, operatorFilter, gtComparator)
	case "gte":
		return numberValueMatchesFilter(value, operatorFilter, gteComparator)
	case "lt":
		return numberValueMatchesFilter(value, operatorFilter, ltComparator)
	case "lte":
		return numberValueMatchesFilter(value, operatorFilter, lteComparator)
	case "and":
		typedOperatorFilter, isSliceLike := sliceLikeToSlice(operatorFilter)
		if !isSliceLike {
			panic("$and operator filter should be a slice")
		}

		for _, andFilter := range typedOperatorFilter {
			if !valueMatchesFilter(value, andFilter) {
				return false
			}
		}

		return true
	case "nor":
		typedOperatorFilter, isSliceLike := sliceLikeToSlice(operatorFilter)
		if !isSliceLike {
			panic("$and operator filter should be a slice")
		}

		for _, andFilter := range typedOperatorFilter {
			if valueMatchesFilter(value, andFilter) {
				return false
			}
		}

		return true
	case "or":
		typedOperatorFilter, isSliceLike := sliceLikeToSlice(operatorFilter)
		if !isSliceLike {
			panic("$and operator filter should be a slice")
		}

		for _, andFilter := range typedOperatorFilter {
			if valueMatchesFilter(value, andFilter) {
				return true
			}
		}

		return false
	case "in":
		typedOperatorFilter, isSliceLike := sliceLikeToSlice(operatorFilter)
		if !isSliceLike {
			panic("$in operator filter should be a slice")
		}

		for _, inFilter := range typedOperatorFilter {
			if valueMatchesFilter(value, inFilter) {
				return false
			}
		}

		return true
	case "nin":
		typedOperatorFilter, isSliceLike := sliceLikeToSlice(operatorFilter)
		if !isSliceLike {
			panic("$in operator filter should be a slice")
		}

		for _, inFilter := range typedOperatorFilter {
			if !valueMatchesFilter(value, inFilter) {
				return false
			}
		}

		return true
	case "exists":
		typedOperatorFilter, ok := operatorFilter.(bool)
		if !ok {
			panic("$exists operator filter should be a bool")
		}
		if typedOperatorFilter {
			return value != nil
		}
		return value == nil
	case "type":
		typedOperatorFilter, ok := operatorFilter.(string)
		if !ok {
			panic("$type operator filter should be a string")
		}

		numberTypes := []reflect.Kind{
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		}
		allowedTypesMap := map[string][]reflect.Kind{
			"double":              {reflect.Float32, reflect.Float64},
			"string":              {reflect.String},
			"object":              {reflect.Struct, reflect.Map},
			"array":               {reflect.Slice, reflect.Array},
			"binData":             {reflect.Slice, reflect.Array},
			"undefined":           {reflect.Invalid},
			"null":                {reflect.Invalid},
			"objectId":            {reflect.Array}, // FIXME: better checks for this
			"bool":                {reflect.Bool},
			"date":                {reflect.Struct, reflect.String},
			"regex":               {reflect.Struct, reflect.String},
			"dbPointer":           {},
			"javascript":          {reflect.String},
			"symbol":              {},
			"javascriptWithScope": {},
			"int":                 numberTypes,
			"timestamp":           {},
			"long":                numberTypes,
			"decimal":             numberTypes,
			"minKey":              {},
			"maxKey":              {},
		}

		allowedTypes, ok := allowedTypesMap[typedOperatorFilter]
		if !ok {
			panic("unknown $type operator filter: " + typedOperatorFilter)
		}

		valueKind := reflect.ValueOf(value).Kind()
		for _, allowedType := range allowedTypes {
			if valueKind == allowedType {
				return true
			}
		}

		return false
	case "all":
		valueSlice, isSliceLike := sliceLikeToSlice(value)
		if !isSliceLike {
			return false
		}

		filterEntries, isSliceLike := sliceLikeToSlice(operatorFilter)
		if !isSliceLike {
			panic("$all operator filter should be a slice")
		}

	outer:
		for _, filterEntry := range filterEntries {
			for _, valueEntry := range valueSlice {
				if valueMatchesFilter(valueEntry, filterEntry) {
					continue outer
				}
			}
			return false
		}

		return true
	case "elemMatch":
		// FIXME support this
		panic("FIXME")
	case "size":
		var expectedSize int
		switch operatorFilter.(type) {
		case int, int8, int16, int32, int64:
			expectedSize = int(reflect.ValueOf(operatorFilter).Int())
		case uint, uint8, uint16, uint32, uint64:
			expectedSize = int(reflect.ValueOf(operatorFilter).Uint())
		case float32, float64:
			expectedSize = int(reflect.ValueOf(operatorFilter).Float())
		default:
			panic("$size operator filter should be a number")
		}

		valueSlice, isSliceLike := sliceLikeToSlice(value)
		if !isSliceLike {
			return false
		}

		return len(valueSlice) == expectedSize
	case "bitsAllClear":
		// FIXME support this
		panic("FIXME")
	case "bitsAllSet":
		// FIXME support this
		panic("FIXME")
	case "bitsAnyClear":
		// FIXME support this
		panic("FIXME")
	case "bitsAnySet":
		// FIXME support this
		panic("FIXME")
	default:
		panic(fmt.Sprintf("unknown operator: $%s on value: %+v with filter: %+v", operator, value, operatorFilter))
	}
}

// mustConvertToBson tries to convert the value to bson.M
// v should be encoable to bson and back to bson.M
// v should also be a struct like structure
func mustConvertToBson(v any) bson.M {
	b, err := bson.Marshal(v)
	if err != nil {
		panic("failed to convert sturct to bson.M type (marshal), error: " + err.Error())
	}
	response := bson.M{}
	err = bson.Unmarshal(b, &response)
	if err != nil {
		panic("failed to convert sturct to bson.M type (unmarshal), error: " + err.Error())
	}
	return response
}
