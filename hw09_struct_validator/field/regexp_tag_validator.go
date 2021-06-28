package field

import (
	"reflect"

	errors2 "github.com/IASamoylov/otus_home_work/hw09_struct_validator/errors"
)

// validateRegexp validates that tag configured for correct type.
func (v validator) validateRegexp(tag Tag) error {
	if v.field.Value.Kind() != reflect.String {
		return errors2.NewValidatorErrorF(
			"tag `%v` not supported for this type %T", tag.Tag, v.field.Value.Interface())
	}

	if tag.Regexp.err != nil {
		return errors2.NewValidatorErrorWF(
			"tag `%v` contains an invalid rule value %v for this type %T",
			tag.Regexp.err, tag.Tag, tag.Value, v.field.Value.Interface())
	}

	return nil
}
