package utils_common

import (
	"fmt"
	"reflect"
)

// Dereference returns the raw value of a variable if it's from a pointer
func Dereference(data interface{}) interface{} {
	if reflect.ValueOf(data).Kind() == reflect.Ptr {
		return reflect.Indirect(reflect.ValueOf(data)).Interface()
	}
	return data
}

// GetStructAsValue returns struct as value
func GetStructAsValue(_struct interface{}) (interface{}, error) {
	rawValue := Dereference(_struct)
	if reflect.TypeOf(rawValue).Kind() != reflect.Struct {
		actualType := reflect.TypeOf(_struct)
		return nil, fmt.Errorf("data provided is not struct: %v, but is of type: %v", _struct, actualType)
	}
	return rawValue, nil
}

// Contains checks if a slice contains a particular element.
func Contains[T comparable](slice []T, compare func(v T) bool) bool {
	for _, v := range slice {
		if compare(v) {
			return true
		}
	}
	return false
}
