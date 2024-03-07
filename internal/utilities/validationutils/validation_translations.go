package validationutils

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

// registerCustomTranslations adds a custom response for the specific validator "name"
func registerCustomTranslations(v *validator.Validate, name, response string) {
	_ = v.RegisterTranslation(name, trans, func(ut ut.Translator) error {
		return ut.Add(name, fmt.Sprintf("{0} %v", response), true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(name, fe.Field())
		return t
	})
}

// registerCustomValidations adds a custom logic
func registerCustomValidations(v *validator.Validate, name string, logic func(i interface{}) bool) {
	_ = v.RegisterValidation(name, func(fl validator.FieldLevel) bool {
		return logic(fl.Field().Interface())
	})
}

// configValidate adds custom validator rules, and custom tag name returns
func configValidate() {
	uni = ut.New(en.New(), en.New())
	trans, _ = uni.GetTranslator("en")
	validate = validator.New()

	setInitialConfiguration(validate)
	addCustomValidations(validate, &customValidations)
}

func addCustomValidations(v *validator.Validate, customValidations *[]customValidation) {
	for _, customValidation := range *customValidations {
		if customValidation.Logic.Valid {
			registerCustomValidations(v, customValidation.Name, customValidation.Logic.Func)
		}
		if customValidation.Response.Valid {
			registerCustomTranslations(v, customValidation.Name, customValidation.Response.String)
		}
	}
}

func setInitialConfiguration(v *validator.Validate) {
	// Display the "json" tag instead of the struct tag  when there are errors
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Set the "required" response error to something more user-friendly
	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "'{0}' must have a value", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
}
