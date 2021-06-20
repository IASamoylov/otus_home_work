package field

import (
	"reflect"
	"strconv"

	"github.com/IASamoylov/otus_home_work/hw09_struct_validator/errors"
)

var validTags = map[string]string{
	InTagValidation:     InTagValidation,
	LenTagValidation:    LenTagValidation,
	MinTagValidation:    MinTagValidation,
	MaxTagValidation:    MaxTagValidation,
	RegexpTagValidation: RegexpTagValidation,
	NestedTagValidation: NestedTagValidation,
}

type validator struct {
	field *Field
}

// newTagValidator ...
func newTagValidator(f *Field) *validator {
	return &validator{field: f}
}

var intConverter = map[reflect.Kind]func(s string) (interface{}, error){
	reflect.Int: func(s string) (interface{}, error) {
		return strconv.Atoi(s)
	},
	reflect.Int8: func(s string) (interface{}, error) {
		return strconv.ParseInt(s, 10, 8)
	},
	reflect.Int16: func(s string) (interface{}, error) {
		return strconv.ParseInt(s, 10, 16)
	},
	reflect.Int32: func(s string) (interface{}, error) {
		return strconv.ParseInt(s, 10, 32)
	},
	reflect.Int64: func(s string) (interface{}, error) {
		return strconv.ParseInt(s, 10, 64)
	},
	reflect.Uint: func(s string) (interface{}, error) {
		return strconv.ParseUint(s, 10, 32)
	},
	reflect.Uint8: func(s string) (interface{}, error) {
		return strconv.ParseUint(s, 10, 8)
	},
	reflect.Uint16: func(s string) (interface{}, error) {
		return strconv.ParseUint(s, 10, 16)
	},
	reflect.Uint32: func(s string) (interface{}, error) {
		return strconv.ParseUint(s, 10, 32)
	},
	reflect.Uint64: func(s string) (interface{}, error) {
		return strconv.ParseUint(s, 10, 64)
	},
}

// Validate validates that validation tags configured correctly
// todo it is best to configure the linter for the current task.
func (v validator) Validate() (ok bool, err error) {
	for _, tag := range v.field.Tags {
		if ok, err = v.baseValidate(tag); !ok {
			return ok, err
		}

		switch tag.Name {
		case InTagValidation:
			err = v.validateIn(tag)
		case LenTagValidation:
			err = v.validateLen(tag)
		case RegexpTagValidation:
			err = v.validateRegexp(tag)
		case MinTagValidation:
			err = v.validateMin(tag)
		case MaxTagValidation:
			err = v.validateMin(tag)
		case NestedTagValidation:
			err = v.validateNested(tag)
		}
	}

	return err == nil, err
}

func (v validator) baseValidate(tag Tag) (ok bool, err error) {
	if _, ok = validTags[tag.Name]; !ok {
		return ok, errors.NewValidatorErrorF(
			"tag `%v` configured incorrectly for field %v", tag.Tag, v.field.FieldType.Name)
	}

	if tag.ValueIsUndefined && tag.Name != NestedTagValidation {
		return false, errors.NewValidatorErrorF(
			"tag `%v` configured incorrectly validation rule for field %v", tag.Tag, v.field.FieldType.Name)
	}

	return true, nil
}

func (v validator) validateTagValueAsInt(tag Tag) error {
	if _, err := intConverter[v.field.Value.Kind()](tag.Value); err != nil {
		return errors.NewValidatorErrorWF(
			"tag `%v` contains an invalid rule value %v for this type %T", err, tag.Tag, tag.Value, v.field.Value.Interface())
	}

	return nil
}
