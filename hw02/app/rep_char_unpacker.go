package app

import (
	"errors"
	"strconv"
	"strings"
)

// RepCharUnpacker is unpacker of repetitive char
type RepCharUnpacker struct {
}

func writeFew(builder *strings.Builder, number *int, rune rune) {
	if *number != 0 {
		for i := 0; i < *number-1; i++ {
			builder.WriteRune(rune)
		}
		*number = 0
	}
}

// Unpack is method of RepCharUnpacker
// example: "a4bc2d5e" => "aaaabccddddde"
func (m *RepCharUnpacker) Unpack(input string) (string, error) {
	result := strings.Builder{}

	var preRune rune
	number := 0
	escape := false

	for pos, curRune := range input {

		if curRune == '\\' {
			writeFew(&result, &number, preRune)

			if escape {
				result.WriteRune(curRune)
				preRune = curRune
				escape = false
			} else {
				escape = true
			}
			continue
		}

		if escape {
			result.WriteRune(curRune)
			preRune = curRune
			escape = false
			continue
		}

		str := string(curRune)
		num, err := strconv.Atoi(str)
		if err != nil {
			writeFew(&result, &number, preRune)

			result.WriteRune(curRune)
			preRune = curRune

			continue
		}

		if pos == 0 {
			return "", errors.New("Number in zero position")
		}
		if num == 0 && number == 0 {
			return "", errors.New("Number is zero")
		}
		number = number*10 + num

	}

	writeFew(&result, &number, preRune)

	if escape {
		return "", errors.New("Empty escape symbol")
	}
	return result.String(), nil
}
