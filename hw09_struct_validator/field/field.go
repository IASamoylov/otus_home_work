package field

import (
	"reflect"
	"regexp"
)

const (
	TagPrefix           string = "validate"
	InTagValidation     string = "in"
	LenTagValidation    string = "len"
	MinTagValidation    string = "min"
	MaxTagValidation    string = "max"
	RegexpTagValidation string = "regexp"
	NestedTagValidation string = "nested"
)

type Field struct {
	Value     reflect.Value
	FieldType reflect.StructField
	Tags      []Tag
}

type TagRegexp struct {
	Regexp *regexp.Regexp
	err    error
}

type Tag struct {
	Tag              reflect.StructTag
	Name             string
	Value            string
	Regexp           TagRegexp
	ValueIsUndefined bool
}

func New(value reflect.Value, index int) *Field {
	fieldType := value.Type().Field(index)

	return &Field{
		Value:     value.Field(index),
		FieldType: fieldType,
		Tags:      newTagParser(fieldType.Tag).Parse(),
	}
}

// HasValidationTags checks that tags was configured for the field.
func (v *Field) HasValidationTags() bool {
	return len(v.Tags) != 0
}

// ValidateTags validates that validation tags configured correctly.
// todo it is best to configure the linter for the current task.
func (v *Field) ValidateTags() (_ bool, err error) {
	return newTagValidator(v).Validate()
}
