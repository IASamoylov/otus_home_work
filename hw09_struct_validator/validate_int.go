package hw09structvalidator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IASamoylov/otus_home_work/hw09_struct_validator/errors"
	"github.com/IASamoylov/otus_home_work/hw09_struct_validator/field"
)

func validateInt(f *field.Field) (ok bool, err error) {
	for _, tag := range f.Tags {
		switch tag.Name {
		case field.MinTagValidation:
			ok, err = validateIntMin(f, tag)
		case field.MaxTagValidation:
			ok, err = validateIntMax(f, tag)
		case field.InTagValidation:
			ok, err = validateIntIn(f, tag)
		}

		if !ok {
			return ok, err
		}
	}

	return true, nil
}

func validateIntMin(field *field.Field, tag field.Tag) (bool, *errors.ValidationError) {
	rule, _ := strconv.ParseInt(tag.Value, 10, 64)
	if field.Value.Int() < rule {
		return false, errors.NewValidationError(
			field.FieldType.Name,
			fmt.Errorf("value must be greater or equal %s", tag.Value))
	}

	return true, nil
}

func validateIntMax(field *field.Field, tag field.Tag) (bool, *errors.ValidationError) {
	rule, _ := strconv.ParseInt(tag.Value, 10, 64)
	if field.Value.Int() > rule {
		return false, errors.NewValidationError(
			field.FieldType.Name,
			fmt.Errorf("value must be les or equal %s", tag.Value))
	}

	return true, nil
}

func validateIntIn(field *field.Field, tag field.Tag) (bool, *errors.ValidationError) {
	value := field.Value.Int()
	r := strings.Split(tag.Value, ",")
	ruleLeft, _ := strconv.ParseInt(r[0], 10, 64)
	ruleRight, _ := strconv.ParseInt(r[1], 10, 64)

	if value < ruleLeft || value > ruleRight {
		return false, errors.NewValidationError(
			field.FieldType.Name,
			fmt.Errorf("value must be in range [%s]", tag.Value))
	}

	return true, nil
}
