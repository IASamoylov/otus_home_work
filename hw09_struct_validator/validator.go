package hw09structvalidator

import (
	"errors"
	"reflect"
)

const (
	tagPrefix           string = "validate"
	inTagValidation     string = "in"
	lenTagValidation    string = "len"
	minTagValidation    string = "min"
	maxTagValidation    string = "max"
	regexpTagValidation string = "regexp"
	nestedTagValidation string = "nested"
)

type ValidationMask string

// Validate validates structure by tag validate:.
func Validate(vi interface{}) error {
	value, err := extractStructure(vi)
	if err != nil {
		return err
	}

	errs := make(ValidationErrors, 0, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		field := newValidationField(value, i)

		if !field.hasValidationTags() {
			continue
		}

		if ok, err := field.validateTags(); !ok {
			return err
		}

		ok, err := validate(field)

		if !ok {
			var validatorErr ValidatorError
			if ok = errors.As(err, &validatorErr); ok {
				return validatorErr
			}

			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

func extractStructure(vi interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(vi)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		return v, nil
	}

	return reflect.Value{}, NewValidatorErrorF("the object %T cannot be validated because it is not a structure", vi)
}

func validate(field *validationField) (bool, *ValidationError) {
	return false, nil
}

func validateStringLen() (bool, *ValidationError) {
	return false, nil
}

func validateStringRegex() (bool, *ValidationError) {
	return false, nil
}

func validateStringIn() (bool, *ValidationError) {
	return false, nil
}

func validateInt32Min() (bool, *ValidationError) {
	return false, nil
}

func validateInt32Max() (bool, *ValidationError) {
	return false, nil
}

func validateInt32In() (bool, *ValidationError) {
	return false, nil
}

func validateInt64Min() (bool, *ValidationError) {
	return false, nil
}

func validateInt64Max() (bool, *ValidationError) {
	return false, nil
}

func validateInt64In() (bool, *ValidationError) {
	return false, nil
}
