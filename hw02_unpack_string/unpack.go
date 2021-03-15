package hw02unpackstring

import (
	"errors"
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

	// Optimizes memory allocation
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

	for _, v := range input {
		switch {
		case unicode.IsDigit(v):
			l, _ := strconv.Atoi(string(v))

			newCapacity += charCap * (l - 1)
			charCap = 0
		default:
			charCap = len(string(v))
			newCapacity += charCap
		}
	}

	return newCapacity
}

func validate(input string) error {
	//  The string cannot start with the digital symbol.
	if len(input) != 0 && unicode.IsDigit([]rune(input)[0]) {
		return ErrInvalidString
	}

	// The string cannot contain numbers.
	if matched, _ := regexp.Match(`\d{2,}`, []byte(input)); matched {
		return ErrInvalidString
	}

	return nil
}

// 日本語
