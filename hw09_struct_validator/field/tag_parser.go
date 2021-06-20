package field

import (
	"reflect"
	"regexp"
	"strings"
)

type parser struct {
	tag reflect.StructTag
}

// newTagParser ...
func newTagParser(tag reflect.StructTag) *parser {
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

		tagName := tagStructure[0]

		var tagRegexp TagRegexp

		if tagName == RegexpTagValidation {
			regexp, err := regexp.Compile(value)
			tagRegexp = TagRegexp{regexp, err}
		}

		fieldTags = append(fieldTags, Tag{
			Tag:              p.tag,
			Name:             tagName,
			Value:            value,
			Regexp:           tagRegexp,
			ValueIsUndefined: valueIsUndefined,
		})
	}

	return fieldTags
}
