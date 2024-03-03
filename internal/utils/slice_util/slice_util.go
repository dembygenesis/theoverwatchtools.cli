package slice_util

import (
	"encoding/json"
	"github.com/friendsofgo/errors"
	"reflect"
)

// Decode decodes the "in", into the "out".
// Basically this is the equivalent, or marshalling something to JSON,
// and unmarshalling it to another struct.
func Decode(in interface{}, out interface{}) error {
	var jsonVal []byte
	var err error

	switch reflect.TypeOf(in).Kind() {
	case reflect.String:
		jsonVal = []byte(in.(string))
	case reflect.Slice:
		// We need to check if the slice actually contains bytes.
		s, ok := in.([]byte)
		if !ok {
			return errors.New("in was a slice but not of bytes")
		}
		jsonVal = s
	default:
		jsonVal, err = json.Marshal(in)
		if err != nil {
			return errors.Wrap(err, "error marshalling 'in'")
		}
	}

	err = json.Unmarshal(jsonVal, out)
	if err != nil {

		return errors.Wrap(err, "error unmarshalling 'out'")
	}
	return nil
}

func Compare[T comparable](slice []T, comparator func(v T) bool) bool {
	for _, v := range slice {
		if comparator(v) {
			return true
		}
	}
	return false
}
