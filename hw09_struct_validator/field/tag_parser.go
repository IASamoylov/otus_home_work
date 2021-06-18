package field

import (
	"reflect"
	"strings"
)

type parser struct {
	tag reflect.StructTag
}

// NewTagParser ...
func NewTagParser(tag reflect.StructTag) *parser {
	return &parser{tag}
}

// Parse converts reflect.StructTag to slice of Tag.
func (p *parser) Parse() []Tag {
	tag, ok := p.tag.Lookup(TagPrefix)

	if !ok {
		return []Tag{}
	}

	tags := strings.Split(tag, "|")

	fieldTags := make([]Tag, 0, len(tags))

	for _, tag := range tags {
		tagStructure := strings.Split(tag, ":")

		var value string
		valueIsUndefined := true
		if len(tagStructure) > 1 && len(tagStructure[1]) > 0 {
			value = tagStructure[1]
			valueIsUndefined = false
		}

		fieldTags = append(fieldTags, Tag{
			Tag:              p.tag,
			Name:             tagStructure[0],
			Value:            value,
			ValueIsUndefined: valueIsUndefined,
		})
	}

	return fieldTags
}
