package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/uchennaemeruche/go-bank-api/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrencry(currency)
	}
	return false
}
