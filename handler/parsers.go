package handler

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

func validateStruct(v *validator.Validate, d interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := v.Struct(d)
	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			var element ErrorResponse
			element.Message = err.Error()
			errors = append(errors, &element)
			return errors
		}
		for _, err := range errs {
			var element ErrorResponse
			element.Field = err.Field()
			element.Value = err.Value()
			element.Message = err.Error()
			errors = append(errors, &element)
		}
	}
	return errors
}
