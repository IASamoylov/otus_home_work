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
	tag              string
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

func (v validationField) hasValidationTags() bool {
	return len(v.tags) != 0
}

// validateTags validates that validation tags configured correctly
// todo it is best to configure the linter for the current task
func (v validationField) validateTags() (bool, error) {
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

		if tag.name == inValidation {
			valuesRange := strings.Split(tag.value, ",")

			if len(valuesRange) != 2 {
				ruleErr := errors.New("rule must be configured as range in:min,max")
				return false, NewValidatorErrorWF("tag %v configured incorrectly validation rule for field %v", ruleErr, tag.tag, v.fieldType.Name)
			}
		}

		kind := v.field.Kind()

		if (tag.name == regexpValidation || tag.name == lenValidation) && kind != reflect.String {
			return false, NewValidatorErrorF("tag %v not supported for this type %T", tag.tag, v.field.Kind())
		}

		if (tag.name == minValidation || tag.name == maxValidation) && kind != reflect.Int32 && kind != reflect.Int64 {
			return false, NewValidatorErrorF("tag %v not supported for this type %T", tag.tag, v.field.Kind())
		}

		if (tag.name == minValidation || tag.name == maxValidation) && (kind == reflect.Int32 || kind == reflect.Int64) {
			if _, err := strconv.Atoi(tag.value); err != nil {
				return false, NewValidatorErrorWF("tag %v contains an invalid rule value %v", err, tag.tag, tag.value)
			}
		}
	}

	return true, nil
}

// parseTags converts reflect.StructTag to slice of validationFieldTag
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
		if len(tagStructure) > 1 {
			value = tagStructure[1]
			valueIsUndefined = false
		}

		validationFieldTags = append(validationFieldTags, validationFieldTag{
			tag:              tag,
			name:             tagStructure[0],
			value:            value,
			valueIsUndefined: valueIsUndefined,
		})
	}

	return validationFieldTags
}
