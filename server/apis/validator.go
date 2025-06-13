package api

import (
	"product-service/utils"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return utils.IsSupportedCurrency(currency)
}
