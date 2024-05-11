package validationutils

import (
	"github.com/dembygenesis/local.tools/internal/utilities/dateutil"
	"github.com/dembygenesis/local.tools/internal/utilities/interfaceutil"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/volatiletech/null/v8"
	"net/url"
	"strings"
	"time"
)

type nullableFunc struct {
	Valid bool
	Func  func(i interface{}) bool
}

type customValidation struct {
	Name     string
	Logic    nullableFunc
	Response null.String
}

var customValidations = []customValidation{
	{
		Name: "is_positive_time_duration",
		Logic: nullableFunc{
			Valid: true,
			Func: func(i interface{}) bool {
				val, ok := i.(time.Duration)
				if !ok {
					return false
				}
				return val > 0
			},
		},
		Response: null.String{
			String: "must be a time duration greater than 0",
			Valid:  true,
		},
	},
	{
		Name: "boolean",
		Logic: nullableFunc{
			Valid: true,
			Func: func(i interface{}) bool {
				_, ok := i.(bool)
				if !ok {
					return false
				}
				return true
			},
		},
		Response: null.String{
			String: "must be a time duration greater than 0",
			Valid:  true,
		},
	},
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
				return strutil.IsValidFilename(i.(string)) == nil
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
				return err == nil
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
				return dateutil.IsValid(val)
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
		Name: "is_positive",
		Logic: nullableFunc{
			Valid: true,
			Func: func(i interface{}) bool {
				nullInt, ok := i.(null.Int)
				if ok {
					if nullInt.Valid {
						return interfaceutil.IsPositive(nullInt.Int)
					}
					return false
				}

				return interfaceutil.IsPositive(i)
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
				return interfaceutil.IsNumericAndGreaterThanZero(i)
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
				return strutil.IsValidFilename(val) == nil
			},
		},
		Response: null.String{
			String: "must be a valid email",
			Valid:  true,
		},
	},
}
