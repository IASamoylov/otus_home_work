package hw02unpackstring

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if err := validate(input); err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.Grow(capacity(input))

	var prevChar string
	for i, v := range input {
		if unicode.IsDigit(v) {
			l, _ := strconv.Atoi(string(v))

			sb.WriteString(strings.Repeat(prevChar, l))

			continue
		}
		if nextIsNotDigit(i, v, input) {
			sb.WriteRune(v)
		}

		prevChar = string(v)
	}

	fmt.Println(sb.Len(), sb.Cap())

	return sb.String(), nil
}

func nextIsNotDigit(i int, v rune, input string) bool {
	l := len(string(v))

	if l+i >= len(input) {
		return true
	}

	return !unicode.IsDigit(rune(input[i+l]))
}

func capacity(input string) int {
	newCapacity := 0
	charCap := 0
	var _ rune

	for _, v := range input {
		switch {
		case unicode.IsDigit(v):
			l, _ := strconv.Atoi(string(v))

			newCapacity += charCap * l
			charCap = 0
		default:
			newCapacity += charCap
			charCap = len(string(v))
		}

		_ = v
	}

	return newCapacity
}

func validate(input string) error {
	if len(input) != 0 && unicode.IsDigit([]rune(input)[0]) {
		return ErrInvalidString
	}

	matched, _ := regexp.Match(`\d{2,}`, []byte(input))

	if matched {
		return ErrInvalidString
	}

	return nil
}

// 日本語
