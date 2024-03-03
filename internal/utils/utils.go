package utils

import "reflect"

// Dereference returns the raw value of a variable if it's from a pointer
func Dereference(data interface{}) interface{} {
	if reflect.ValueOf(data).Kind() == reflect.Ptr {
		return reflect.Indirect(reflect.ValueOf(data)).Interface()
	}
	return data
}
