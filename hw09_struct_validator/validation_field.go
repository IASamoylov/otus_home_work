package hw09structvalidator

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type validationField struct {
	field     reflect.Value
	fieldType reflect.StructField
	tags      []validationFieldTag
}

type validationFieldTag struct {
	tag              reflect.StructTag
	name             string
	value            string
	valueIsUndefined bool
}

func newValidationField(value reflect.Value, index int) *validationField {
	fieldType := value.Type().Field(index)

	return &validationField{
		field:     value.Field(index),
		fieldType: fieldType,
		tags:      parseTags(fieldType.Tag),
	}
}

func (v *validationField) hasValidationTags() bool {
	return len(v.tags) != 0
}

// validateTags validates that validation tags configured correctly
// todo it is best to configure the linter for the current task.
func (v *validationField) validateTags() (_ bool, err error) {
	validTags := map[string]string{
		inValidation:     inValidation,
		lenValidation:    lenValidation,
		minValidation:    minValidation,
		maxValidation:    maxValidation,
		regexpValidation: regexpValidation,
	}

	for _, tag := range v.tags {
		if _, ok := validTags[tag.name]; !ok {
			return false, NewValidatorErrorF("tag %v configured incorrectly for field %v", tag.tag, v.fieldType.Name)
		}

		if tag.valueIsUndefined {
			return false, NewValidatorErrorF("tag %v configured incorrectly validation rule for field %v", tag.tag, v.fieldType.Name)
		}

		switch tag.name {
		case inValidation:
			err = validateIn(v, tag)
		case lenValidation:
			err = validateLen(v, tag)
		case regexpValidation:
			err = validateRegexp(v, tag)
		case minValidation:
			err = validateMinmax(v, tag)
		case maxValidation:
			err = validateMinmax(v, tag)
		}
	}

	return true, err
}

// validateIn validates that tag configured for correct type.
func validateIn(v *validationField, tag validationFieldTag) error {
	valuesRange := strings.Split(tag.value, ",")

	if len(valuesRange) != 2 {
		ruleErr := errors.New("rule must be configured as range in:min,max")
		return NewValidatorErrorWF("tag %v configured incorrectly validation rule for field %v", ruleErr, tag.tag, v.fieldType.Name)
	}

	return nil
}

// validateLen validates that tag configured for correct type.
func validateLen(v *validationField, tag validationFieldTag) error {
	if v.field.Kind() != reflect.String {
		return NewValidatorErrorF("tag %v not supported for this type %T", tag.tag, v.field.Kind())
	}
	return nil
}

// validateRegexp validates that tag configured for correct type.
func validateRegexp(v *validationField, tag validationFieldTag) error {
	if v.field.Kind() != reflect.String {
		return NewValidatorErrorF("tag %v not supported for this type %T", tag.tag, v.field.Kind())
	}
	return nil
}

// validateMinmax validates that tag configured for correct type.
func validateMinmax(v *validationField, tag validationFieldTag) error {
	kind := v.field.Kind()
	if kind != reflect.Int8 && kind != reflect.Int16 && kind != reflect.Int32 && kind != reflect.Int64 &&
		kind != reflect.Uint8 && kind != reflect.Uint16 && kind != reflect.Uint32 && kind != reflect.Uint64 {
		return NewValidatorErrorF("tag %v not supported for this type %T", tag.tag, v.field.Kind())
	}

	types := map[reflect.Kind]func(s string) (interface{}, error){
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
	for k, parse := range types {
		if _, err := parse(tag.value); kind == k && err != nil {
			return NewValidatorErrorWF("tag %v contains an invalid rule value %v", err, tag.tag, tag.value)
		}
	}

	return nil
}

// parseTags converts reflect.StructTag to slice of validationFieldTag.
func parseTags(structTag reflect.StructTag) []validationFieldTag {
	tag, ok := structTag.Lookup(tagPrefix)

	if !ok {
		return []validationFieldTag{}
	}

	tags := strings.Split(tag, "|")

	validationFieldTags := make([]validationFieldTag, 0, len(tags))

	for _, tag := range tags {
		tagStructure := strings.Split(tag, ":")

		var value string
		valueIsUndefined := true
		if len(tagStructure) > 1 && len(tagStructure[1]) > 0 {
			value = tagStructure[1]
			valueIsUndefined = false
		}

		validationFieldTags = append(validationFieldTags, validationFieldTag{
			tag:              structTag,
			name:             tagStructure[0],
			value:            value,
			valueIsUndefined: valueIsUndefined,
		})
	}

	return validationFieldTags
}
