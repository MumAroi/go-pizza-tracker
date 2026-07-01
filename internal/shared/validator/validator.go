package validators

import (
	"pizza-tracker/internal/order"
	"slices"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
)

func RegisterCustomValidators() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}

	v.RegisterValidation("valid_pizza_type", createSliceValidator(order.GetPizzaTypes()))
	v.RegisterValidation("valid_pizza_size", createSliceValidator(order.GetPizzaSizes()))
}

func createSliceValidator(allowedValues []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return slices.Contains(allowedValues, fl.Field().String())
	}
}
