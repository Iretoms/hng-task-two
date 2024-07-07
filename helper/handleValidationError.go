package helper

import (
	"github.com/Iretoms/hng-task-two/responses"
	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) []responses.Error {
	var errors []responses.Error
	for _, err := range err.(validator.ValidationErrors) {
		var element responses.Error
		element.Field = err.Field()
		element.Message = "Invalid input"
		errors = append(errors, element)
	}
	return errors
}
