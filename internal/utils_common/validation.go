package utils_common

import (
	"errors"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func ValidateStruct(p interface{}, exclusions ...string) error {
	var errMsgs ErrorList

	structVal, err := GetStructAsValue(p)
	if err != nil {
		errMsgs.AddErr(err)
		return errMsgs.Single()
	}

	err = validate.Struct(structVal)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, err := range validationErrors {
			translatedErr := err.Translate(trans)
			if strings.Contains(translatedErr, "__VAL__") {
				val, ok := err.Value().(string)
				if ok {
					translatedErr = strings.ReplaceAll(translatedErr, "__VAL__", val)
				}
			}
			if translatedErr != "" {
				errMsgs = append(errMsgs, errors.New(translatedErr))
			} else {
				errMsgs = append(errMsgs, err)
			}
		}
	}

	if len(exclusions) > 0 {
		filtered := make([]string, 0)
		excludeHash := makeHashVariadic(exclusions...)
		for _, errMsg := range errMsgs {
			if !excludeHash[errMsg] {
				filtered = append(filtered, errMsg.Error())
			}
		}
	}

	return errMsgs.Single()
}

// makeHashVariadic removes an element specified by index for any type of slice
func makeHashVariadic[T any](variadic ...T) map[any]bool {
	arr := make([]T, 0)
	for _, val := range variadic {
		arr = append(arr, val)
	}
	return makeHash(arr)
}

// makeHash removes an element specified by index for any type of slice
func makeHash[T any](s []T) map[any]bool {
	m := make(map[any]bool)
	for _, v := range s {
		m[v] = true
	}

	return m
}
