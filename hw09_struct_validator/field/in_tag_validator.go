package field

import (
	"errors"
	"reflect"
	"strings"

	validatorerrors "github.com/IASamoylov/otus_home_work/hw09_struct_validator/errors"
)

// validateIn validates that tag configured for correct type.
func (v validator) validateIn(tag Tag) error {
	switch v.field.Value.Kind() {
	case reflect.String:
		return v.validateInString(tag)
	case reflect.Int:
		return v.validateInInt(tag)
	case reflect.Int8:
		return v.validateInInt(tag)
	case reflect.Int16:
		return v.validateInInt(tag)
	case reflect.Int32:
		return v.validateInInt(tag)
	case reflect.Int64:
		return v.validateInInt(tag)
	case reflect.Uint:
		return v.validateInInt(tag)
	case reflect.Uint8:
		return v.validateInInt(tag)
	case reflect.Uint16:
		return v.validateInInt(tag)
	case reflect.Uint32:
		return v.validateInInt(tag)
	case reflect.Uint64:
		return v.validateInInt(tag)
	}

	return validatorerrors.NewValidatorErrorF(
		"tag `%v` not supported for this type %T", tag.Tag, v.field.Value.Interface())
}

func (v validator) validateInString(tag Tag) error {
	values := strings.Split(tag.Value, ",")

	if len(values) == 0 {
		ruleErr := errors.New("rule must be configured as range in:value1,value2,value3")
		return validatorerrors.NewValidatorErrorWF(
			"tag `%v` configured incorrectly validation rule for field %s", ruleErr, tag.Tag, v.field.FieldType.Name)
	}

	return nil
}

func (v validator) validateInInt(tag Tag) error {
	values := strings.Split(tag.Value, ",")

	if len(values) != 2 || len(values[1]) == 0 {
		ruleErr := errors.New("rule must be configured as range in:min,max")
		return validatorerrors.NewValidatorErrorWF(
			"tag `%v` configured incorrectly validation rule for field %s", ruleErr, tag.Tag, v.field.FieldType.Name)
	}

	for _, value := range values {
		if _, err := intConverter[v.field.Value.Kind()](value); err != nil {
			return validatorerrors.NewValidatorErrorWF(
				"tag `%v` contains an invalid rule value %v for this type %T", err, tag.Tag, tag.Value, v.field.Value.Interface())
		}
	}

	return nil
}
