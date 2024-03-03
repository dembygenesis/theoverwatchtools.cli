package interface_util

import (
	"fmt"
	"reflect"
)

// IsNumericAndGreaterThanZero checks if the input is a numeric primitive and greater than zero.
// Returns true if the condition is met, otherwise false.
func IsNumericAndGreaterThanZero(i interface{}) bool {
	switch v := i.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() > 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint() > 0
	case float32, float64:
		return reflect.ValueOf(v).Float() > 0
	default:
		return false
	}
}

// Dereference returns the raw value of a variable if it's from a pointer
func Dereference(i interface{}) interface{} {
	if reflect.ValueOf(i).Kind() == reflect.Ptr {
		return reflect.Indirect(reflect.ValueOf(i)).Interface()
	}
	return i
}

// GetStructAsValue returns a struct's raw value.
func GetStructAsValue(i interface{}) (interface{}, error) {
	val := Dereference(i)
	if reflect.TypeOf(val).Kind() != reflect.Struct {
		actualType := reflect.TypeOf(i)
		return nil, fmt.Errorf("data provided is not struct: %v, but is of type: %v", i, actualType)
	}
	return val, nil
}
