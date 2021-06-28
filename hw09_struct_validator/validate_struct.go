package hw09structvalidator

import (
	"errors"

	errors2 "github.com/IASamoylov/otus_home_work/hw09_struct_validator/errors"
	"github.com/IASamoylov/otus_home_work/hw09_struct_validator/field"
)

func validateStruct(field *field.Field) (bool, error) {
	if ok, err := validate(field.Value.Interface()); !ok {
		var validatorErr *errors2.ValidatorError
		if ok = errors.As(err, &validatorErr); ok {
			return false, validatorErr
		}
		return false, errors2.NewValidationError(field.FieldType.Name, err)
	}

	return true, nil
}
