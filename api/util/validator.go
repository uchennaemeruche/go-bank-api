package api

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/uchennaemeruche/go-bank-api/util"
)

var ValidCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrencry(currency)
	}
	return false
}

type ApiError struct {
	Msg      string
	Expected string
	Got      string
}

func FormatValidationErr(errs validator.ValidationErrors) []ApiError {

	result := make([]ApiError, len(errs))

	for i, fieldErr := range errs {
		var expected strings.Builder
		expected.WriteString(fieldErr.ActualTag())
		if fieldErr.Param() != "" {
			expected.WriteString(" { " + fieldErr.Param() + " }")
		}

		var actual strings.Builder
		if fieldErr.Value() != nil && fieldErr.Value() != "" {
			actual.WriteString(fmt.Sprintf(" %v", fieldErr.Value()))
		}

		result[i] = ApiError{
			Msg:      "Input validation failed on field '" + fieldErr.Field() + "'",
			Expected: expected.String(),
			Got:      actual.String(),
		}

	}
	return result
}
