package hw09structvalidator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IASamoylov/otus_home_work/hw09_struct_validator/errors"
	"github.com/IASamoylov/otus_home_work/hw09_struct_validator/field"
)

func validateString(f *field.Field) (ok bool, err error) {
	for _, tag := range f.Tags {
		switch tag.Name {
		case field.LenTagValidation:
			return validateStringLen(f, tag)
		case field.RegexpTagValidation:
			return validateStringRegex(f, tag)
		case field.InTagValidation:
			return validateStringIn(f, tag)
		}

		if !ok {
			return ok, err
		}
	}

	return true, nil
}

func validateStringLen(field *field.Field, tag field.Tag) (bool, error) {
	value := field.Value.String()
	length, _ := strconv.Atoi(tag.Value)

	if len(value) <= length {
		return false, errors.NewValidationError(
			field.FieldType.Name,
			fmt.Errorf("string length must be greater or equal %s", tag.Value))
	}

	return true, nil
}

func validateStringRegex(field *field.Field, tag field.Tag) (bool, error) {
	if !tag.Regexp.Regexp.MatchString(field.Value.String()) {
		return false, errors.NewValidationError(
			field.FieldType.Name,
			fmt.Errorf("value does not match the mask %s", tag.Value))
	}

	return true, nil
}

func validateStringIn(field *field.Field, tag field.Tag) (bool, error) {
	values := strings.Split(tag.Value, ",")

	m := make(map[string]string, len(values))

	for _, v := range values {
		m[v] = v
	}

	if _, ok := m[field.Value.String()]; !ok {
		return false, errors.NewValidationError(
			field.FieldType.Name,
			fmt.Errorf("value is not contains in range %s", tag.Value))
	}

	return true, nil
}
