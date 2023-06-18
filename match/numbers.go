package match

import (
	"reflect"
)

type numberComparatorT struct {
	Int   func(a, b int64) bool
	Uint  func(a, b uint64) bool
	Float func(a, b float64) bool
}

var eqComparator = numberComparatorT{
	Int:   func(a, b int64) bool { return a == b },
	Uint:  func(a, b uint64) bool { return a == b },
	Float: func(a, b float64) bool { return a == b },
}
var gtComparator = numberComparatorT{
	Int:   func(a, b int64) bool { return a > b },
	Uint:  func(a, b uint64) bool { return a > b },
	Float: func(a, b float64) bool { return a > b },
}
var gteComparator = numberComparatorT{
	Int:   func(a, b int64) bool { return a >= b },
	Uint:  func(a, b uint64) bool { return a >= b },
	Float: func(a, b float64) bool { return a >= b },
}
var ltComparator = numberComparatorT{
	Int:   func(a, b int64) bool { return a < b },
	Uint:  func(a, b uint64) bool { return a < b },
	Float: func(a, b float64) bool { return a < b },
}
var lteComparator = numberComparatorT{
	Int:   func(a, b int64) bool { return a <= b },
	Uint:  func(a, b uint64) bool { return a <= b },
	Float: func(a, b float64) bool { return a <= b },
}

func numberValueMatchesFilter(value any, filter any, numberComparator numberComparatorT) bool {
	switch typedFilter := filter.(type) {
	case int, int8, int16, int32, int64:
		filterInt := reflect.ValueOf(filter).Int()
		return intFilterMatches(value, filterInt, numberComparator)
	case uint, uint8, uint16, uint32, uint64:
		filterUint := reflect.ValueOf(filter).Uint()
		return uintFilterMatches(value, filterUint, numberComparator)
	case float32:
		return floatFilterMatches(value, float64(typedFilter), numberComparator)
	case float64:
		return floatFilterMatches(value, typedFilter, numberComparator)
	default:
		return false
	}
}

func intFilterMatches(value any, filter int64, numberComparator numberComparatorT) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return numberComparator.Int(reflect.ValueOf(value).Int(), filter)
	case uint, uint8, uint16, uint32, uint64:
		return numberComparator.Int(int64(reflect.ValueOf(value).Uint()), filter)
	case float32, float64:
		typedValue := reflect.ValueOf(value).Float()
		if int64(typedValue*1000)%1000 == 0 {
			return numberComparator.Int(int64(typedValue), filter)
		}
		return numberComparator.Float(typedValue, float64(filter))
	default:
		return false
	}
}

func uintFilterMatches(value any, filter uint64, numberComparator numberComparatorT) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		typedValue := reflect.ValueOf(value).Int()
		if typedValue < 0 {
			return false
		}
		return numberComparator.Uint(uint64(typedValue), filter)
	case uint, uint8, uint16, uint32, uint64:
		return numberComparator.Uint(reflect.ValueOf(value).Uint(), filter)
	case float32, float64:
		typedValue := reflect.ValueOf(value).Float()
		if typedValue < 0 {
			return false
		}
		if uint64(typedValue*1000)%1000 == 0 {
			return numberComparator.Uint(uint64(typedValue), filter)
		}
		return numberComparator.Float(typedValue, float64(filter))
	default:
		return false
	}
}

func floatFilterMatches(value any, filter float64, numberComparator numberComparatorT) bool {
	if int64(filter*1000)%1000 == 0 {
		return intFilterMatches(value, int64(filter), numberComparator)
	}

	switch value.(type) {
	case float32, float64:
		return numberComparator.Float(reflect.ValueOf(value).Float(), filter)
	case int, int8, int16, int32, int64:
		return numberComparator.Float(float64(reflect.ValueOf(value).Int()), filter)
	case uint, uint8, uint16, uint32, uint64:
		return numberComparator.Float(float64(reflect.ValueOf(value).Uint()), filter)
	default:
		return false
	}
}
