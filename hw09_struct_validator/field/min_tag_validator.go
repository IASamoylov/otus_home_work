package field

import (
	"reflect"

	errors2 "github.com/IASamoylov/otus_home_work/hw09_struct_validator/errors"
)

// validateMin validates that tag configured for correct type.
func (v validator) validateMin(tag Tag) error {
	switch v.field.Value.Kind() {
	case reflect.Int:
		return v.validateTagValueAsInt(tag)
	case reflect.Int8:
		return v.validateTagValueAsInt(tag)
	case reflect.Int16:
		return v.validateTagValueAsInt(tag)
	case reflect.Int32:
		return v.validateTagValueAsInt(tag)
	case reflect.Int64:
		return v.validateTagValueAsInt(tag)
	case reflect.Uint:
		return v.validateTagValueAsInt(tag)
	case reflect.Uint8:
		return v.validateTagValueAsInt(tag)
	case reflect.Uint16:
		return v.validateTagValueAsInt(tag)
	case reflect.Uint32:
		return v.validateTagValueAsInt(tag)
	case reflect.Uint64:
		return v.validateTagValueAsInt(tag)
	}

	return errors2.NewValidatorErrorF(
		"tag `%v` not supported for this type %T", tag.Tag, v.field.Value.Interface())
}
