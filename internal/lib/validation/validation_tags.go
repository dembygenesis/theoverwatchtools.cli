package validation

import (
	"github.com/dembygenesis/local.tools/internal/utils/date_util"
	"github.com/dembygenesis/local.tools/internal/utils/interface_util"
	"github.com/dembygenesis/local.tools/internal/utils/string_util"
	"github.com/volatiletech/null"
	"net/url"
	"strings"
)

type nullableFunc struct {
	Valid bool
	Func  func(i interface{}) bool
}

type CustomValidation struct {
	Name     string
	Logic    nullableFunc
	Response null.String
}

var customValidations = []CustomValidation{
	{
		Name: "ascending_non_negative_ints",
		Logic: nullableFunc{
			Valid: true,
			Func: func(i interface{}) bool {
				ints, ok := i.([]int)
				if !ok {
					return false
				}
				if len(ints) == 0 {
					return true
				}
				for i := 1; i < len(ints); i++ {
					if ints[i] < ints[i-1] || ints[i] < 0 {
						return false
					}
				}
				return true
			},
		},
		Response: null.String{
			Valid:  true,
			String: "must be an array of non-negative, ascending integers",
		},
	},
	{
		Name: "max",
		Logic: nullableFunc{
			Valid: false,
		},
		Response: null.String{
			Valid:  true,
			String: "string exceeds the max number of characters",
		},
	},
	{
		Name: "is_file_name",
		Logic: nullableFunc{
			Valid: true,
			Func: func(i interface{}) bool {
				return string_util.IsValidFilename(i.(string)) == nil
			},
		},
		Response: null.String{
			Valid: true,
			// Might need to change this template into something else.
			String: "'__VAL__' is not a valid file_name",
		},
	},

	{
		Name: "is_url",
		Logic: nullableFunc{
			Valid: true,
			Func: func(i interface{}) bool {
				val, ok := i.(string)
				if !ok {
					return false
				}
				val = strings.TrimSpace(val)
				if val == "" {
					return false
				}
				_, err := url.ParseRequestURI(val)
				return err != nil
			},
		},
		Response: null.String{
			Valid:  true,
			String: "must be a valid url format",
		},
	},
	{
		Name: "date_format",
		Logic: nullableFunc{
			Valid: true,
			Func: func(i interface{}) bool {
				val, ok := i.(string)
				if !ok {
					return false
				}
				return date_util.IsValidDate(val)
			},
		},
		Response: null.String{
			String: "must be a valid date format of YYYY-MM-DD",
			Valid:  true,
		},
	},
	{
		Name: "is_valid_filename",
		Logic: nullableFunc{
			Valid: true,
			Func: func(i interface{}) bool {
				val, ok := i.(int)
				if !ok {
					return false
				}
				return val > 0
			},
		},
		Response: null.String{
			String: "must be an integer greater than 0",
			Valid:  true,
		},
	},
	{
		Name: "greater_than_zero",
		Logic: nullableFunc{
			Valid: true,
			Func: func(i interface{}) bool {
				return interface_util.IsNumericAndGreaterThanZero(i)
			},
		},
		Response: null.String{
			String: "must be an integer greater than 0",
			Valid:  true,
		},
	},
	{
		Name: "float64_greater_than_zero",
		Logic: nullableFunc{
			Valid: true,
			Func: func(i interface{}) bool {
				val, ok := i.(float64)
				if !ok {
					return false
				}
				return val > 0
			},
		},
		Response: null.String{
			String: "must be a a float64 greater than 0",
			Valid:  true,
		},
	},
	{
		Name: "email",
		Logic: nullableFunc{
			Valid: true,
			Func: func(i interface{}) bool {
				val, ok := i.(string)
				if !ok {
					return false
				}
				return string_util.IsValidFilename(val) == nil
			},
		},
		Response: null.String{
			String: "must be a valid email",
			Valid:  true,
		},
	},
}
