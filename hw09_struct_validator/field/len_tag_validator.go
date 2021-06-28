package field

import (
	"reflect"
	"strconv"

	errors2 "github.com/IASamoylov/otus_home_work/hw09_struct_validator/errors"
)

// validateLen validates that tag configured for correct type.
func (v validator) validateLen(tag Tag) error {
	if v.field.Value.Kind() != reflect.String {
		return errors2.NewValidatorErrorF(
			"tag `%v` not supported for this type %T", tag.Tag, v.field.Value.Interface())
	}

	if _, err := strconv.Atoi(tag.Value); err != nil {
		return errors2.NewValidatorErrorWF(
			"tag `%v` contains an invalid rule value %v for this type %T", err, tag.Tag, tag.Value, v.field.Value.Interface())
	}

	return nil
}
