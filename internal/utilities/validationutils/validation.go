package validationutils

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utilities/errutil"
	"github.com/dembygenesis/local.tools/internal/utilities/interfaceutil"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

var (
	ErrStructNil = errors.New("error, struct is nil")
)

func init() {
	configValidate()
}

func Validate(p interface{}, exclusions ...string) error {
	if p != nil && fmt.Sprintf("%v", p) == "<nil>" {
		return ErrStructNil
	}

	var errs errutil.List

	structVal, err := interfaceutil.GetStructAsValue(p)
	if err != nil {
		errs.AddErr(err)
		return errs.Single()
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
				errs = append(errs, errors.New(translatedErr))
			} else {
				errs = append(errs, err)
			}
		}
	}

	return errs.Single()
}

// makeHashVariadic removes an element specified by index for any type of slice
func makeHashVariadic[T any](variadic ...T) map[any]bool {
	return makeHash(append(make([]T, 0), variadic...))
}

// makeHash removes an element specified by index for any type of slice
func makeHash[T any](s []T) map[any]bool {
	m := make(map[any]bool)
	for _, v := range s {
		m[v] = true
	}

	return m
}
