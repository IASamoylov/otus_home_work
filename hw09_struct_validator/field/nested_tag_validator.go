package field

import (
	"reflect"

	errors2 "github.com/IASamoylov/otus_home_work/hw09_struct_validator/errors"
)

// validateNested validates that tag configured for correct type.
func (v validator) validateNested(tag Tag) error {
	if v.field.Value.Kind() != reflect.Struct && v.field.Value.Kind() != reflect.Ptr && v.field.Value.Kind() != reflect.Interface {
		return errors2.NewValidatorErrorF(
			"tag `%v` not supported for this type %T", tag.Tag, v.field.Value.Interface())
	}
	return nil
}
