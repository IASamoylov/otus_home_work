package hw09structvalidator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IASamoylov/otus_home_work/hw09_struct_validator/errors"
	"github.com/IASamoylov/otus_home_work/hw09_struct_validator/field"
)

func validateUint(f *field.Field) (ok bool, err error) {
	for _, tag := range f.Tags {
		switch tag.Name {
		case field.MinTagValidation:
			ok, err = validateUintMin(f, tag)
		case field.MaxTagValidation:
			ok, err = validateUintMax(f, tag)
		case field.InTagValidation:
			ok, err = validateUintIn(f, tag)
		}

		if !ok {
			return ok, err
		}
	}

	return true, nil
}

func validateUintMin(field *field.Field, tag field.Tag) (bool, *errors.ValidationError) {
	rule, _ := strconv.ParseUint(tag.Value, 10, 64)
	if field.Value.Uint() < rule {
		return false, errors.NewValidationError(
			field.FieldType.Name,
			fmt.Errorf("value must be greater or equal %s", tag.Value))
	}

	return true, nil
}

func validateUintMax(field *field.Field, tag field.Tag) (bool, *errors.ValidationError) {
	rule, _ := strconv.ParseUint(tag.Value, 10, 64)
	if field.Value.Uint() > rule {
		return false, errors.NewValidationError(
			field.FieldType.Name,
			fmt.Errorf("value must be les or equal %s", tag.Value))
	}

	return true, nil
}

func validateUintIn(field *field.Field, tag field.Tag) (bool, *errors.ValidationError) {
	value := field.Value.Uint()
	r := strings.Split(tag.Value, ",")
	ruleLeft, _ := strconv.ParseUint(r[0], 10, 64)
	ruleRight, _ := strconv.ParseUint(r[1], 10, 64)

	if value < ruleLeft || value > ruleRight {
		return false, errors.NewValidationError(
			field.FieldType.Name,
			fmt.Errorf("value must be in range [%s]", tag.Value))
	}

	return true, nil
}
