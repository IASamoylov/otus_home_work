package hw09structvalidator

import (
	"errors"
	"reflect"

	errors2 "github.com/IASamoylov/otus_home_work/hw09_struct_validator/errors"
	"github.com/IASamoylov/otus_home_work/hw09_struct_validator/field"
)

type ValidationMask string

// Validate validates structure by tag validate:.
func Validate(vi interface{}) error {
	if ok, err := validate(vi); !ok {
		return err
	}

	return nil
}

func validate(vi interface{}) (bool, error) {
	value, err := extractStructure(vi)
	if err != nil {
		return false, err
	}

	errs := make([]*errors2.ValidationError, 0, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		field := field.New(value, i)

		if !field.HasValidationTags() {
			continue
		}

		if ok, err := field.ValidateTags(); !ok {
			return false, err
		}

		ok, err := validateField(field)

		if !ok {
			var validatorErr *errors2.ValidatorError
			if ok = errors.As(err, &validatorErr); ok {
				return false, validatorErr
			}
			var validationErr *errors2.ValidationError
			if ok = errors.As(err, &validationErr); ok {
				errs = append(errs, validationErr)
			}

			return false, &errors2.ValidationErrors{Errors: errs}
		}
	}

	if len(errs) != 0 {
		return false, &errors2.ValidationErrors{Errors: errs}
	}

	return true, nil
}

func extractStructure(vi interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(vi)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		return v, nil
	}

	return reflect.Value{}, errors2.NewValidatorErrorF("the object %T cannot be validated because it is not a structure", vi)
}

func validateField(field *field.Field) (bool, error) {
	switch field.Value.Kind() {
	case reflect.Struct:
		return validateStruct(field)
	case reflect.Ptr:
		return validateStruct(field)
	case reflect.String:
		return validateString(field)
	case reflect.Int:
		return validateInt(field)
	case reflect.Int8:
		return validateInt(field)
	case reflect.Int16:
		return validateInt(field)
	case reflect.Int32:
		return validateInt(field)
	case reflect.Int64:
		return validateInt(field)
	case reflect.Uint:
		return validateUint(field)
	case reflect.Uint8:
		return validateUint(field)
	case reflect.Uint16:
		return validateUint(field)
	case reflect.Uint32:
		return validateUint(field)
	case reflect.Uint64:
		return validateUint(field)
	}

	return false, nil
}
